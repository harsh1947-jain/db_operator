package controller

import (
	"context"
	demov1 "github.com/harsh1947-jain/db-operator/api/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type DboperatorReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

func (r *DboperatorReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	l := log.FromContext(ctx)
	db := &demov1.Dboperator{}
	if err := r.Get(ctx, req.NamespacedName, db); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	stsName := db.Spec.Username + "-" + db.Spec.DbType + "-" + db.Spec.DbName
	found := &appsv1.StatefulSet{}
	err := r.Get(ctx, types.NamespacedName{Name: stsName, Namespace: req.Namespace}, found)

	if err != nil && errors.IsNotFound(err) {
		l.Info("Creating StatefulSet", "Name", stsName)
		sts := r.desiredStatefulSet(db, stsName)
		if err := ctrl.SetControllerReference(db, sts, r.Scheme); err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, r.Create(ctx, sts)
	}
	return ctrl.Result{}, err
}

func (r *DboperatorReconciler) desiredStatefulSet(db *demov1.Dboperator, name string) *appsv1.StatefulSet {
	config := getDbDefaults(db)
	replicas := int32(1)
	labels := map[string]string{"app": name}

	return &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: db.Namespace},
		Spec: appsv1.StatefulSetSpec{
			Replicas:    &replicas,
			ServiceName: name,
			Selector:    &metav1.LabelSelector{MatchLabels: labels},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: labels},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Name:  config.NameofDB,
						Image: config.ImageofDB,
						Resources: corev1.ResourceRequirements{
							Requests: corev1.ResourceList{corev1.ResourceCPU: config.CPU, corev1.ResourceMemory: config.Memory},
							Limits:   corev1.ResourceList{corev1.ResourceCPU: config.CPU, corev1.ResourceMemory: config.Memory},
						},
						Env: []corev1.EnvVar{
							{Name: config.UsernameEnv, Value: db.Spec.Username},
							{Name: config.PasswordEnv, Value: db.Spec.Password},
							{Name: config.DbEnv, Value: db.Spec.DbName},
						},
						VolumeMounts: []corev1.VolumeMount{{Name: "data", MountPath: config.Mountingpath}},
					}},
				},
			},
			VolumeClaimTemplates: []corev1.PersistentVolumeClaim{{
				ObjectMeta: metav1.ObjectMeta{Name: "data"},
				Spec: corev1.PersistentVolumeClaimSpec{
					AccessModes: []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
					Resources: corev1.ResourceRequirements{
						Requests: corev1.ResourceList{corev1.ResourceStorage: config.Storage},
					},
				},
			}},
		},
	}
}

func (r *DboperatorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&demov1.Dboperator{}).
		Owns(&appsv1.StatefulSet{}).
		Complete(r)
}
