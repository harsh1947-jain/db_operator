package controller

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1" // Fixes undefined: metav1
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	demov1 "github.com/harsh1947-jain/db-operator/api/v1"
)

type DboperatorReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

func (r *DboperatorReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	// Assigning to 'l' avoids conflict with the 'log' package name
	l := log.FromContext(ctx)

	dboperator := &demov1.Dboperator{}
	err := r.Get(ctx, req.NamespacedName, dboperator)
	if err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	dbtype := dboperator.Spec.DbType

	stsName := dboperator.Spec.Username + "-" + dbtype + "-" + dboperator.Spec.DbName

	found := &appsv1.StatefulSet{}
	err = r.Get(ctx, types.NamespacedName{Name: stsName, Namespace: req.Namespace}, found)

	if err != nil && errors.IsNotFound(err) {
		l.Info("Creating a new StatefulSet", "Namespace", req.Namespace, "Name", stsName)
		newSts := r.desiredStatefulSet(dboperator, stsName)

		if err := ctrl.SetControllerReference(dboperator, newSts, r.Scheme); err != nil {
			return ctrl.Result{}, err
		}

		err = r.Create(ctx, newSts)
		if err != nil {
			l.Error(err, "Failed to create StatefulSet") // Use 'l' instead of 'log'
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
	}

	return ctrl.Result{}, err
}

func (r *DboperatorReconciler) desiredStatefulSet(db *demov1.Dboperator, name string) *appsv1.StatefulSet {
	
	config := getDbDefaults(db)
	replicas := int32(2)
	labels := map[string]string{"app": name}
	storage := resource.MustParse(db.Spec.Storage)

	return &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: db.Namespace,
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas:    &replicas,
			ServiceName: name, // Fixed: ServiceName must be inside Spec
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{ // Added missing curly brace for the container object
							Name:  config.NameofDB,
							Image: config.ImageofDB,
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      config.NameofDB,
									MountPath: config.Mountingpath,
								},
							},
							Env: []corev1.EnvVar{
								{Name: config.usernameenv, Value: db.Spec.Username},
								{Name: config.passwordenv, Value: db.Spec.Password},
								{Name: config.dbenv, Value: db.Spec.DbName},
							},
						},
					},
				},
			},
			VolumeClaimTemplates: []corev1.PersistentVolumeClaim{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: config.NameofDB,
					},
					Spec: corev1.PersistentVolumeClaimSpec{
						AccessModes: []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
						Resources: corev1.ResourceRequirements{
							Requests: corev1.ResourceList{
								corev1.ResourceStorage: storage,
							},
						},
					},
				},
			},
		},
	}
}

func (r *DboperatorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&demov1.Dboperator{}).
		Owns(&appsv1.StatefulSet{}).
		Named("db").
		Complete(r)
}
