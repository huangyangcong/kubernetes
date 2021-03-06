/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package fuzzer

import (
	fuzz "github.com/google/gofuzz"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtimeserializer "k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm"
	"k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm/v1beta1"
	"k8s.io/kubernetes/cmd/kubeadm/app/constants"
)

// Funcs returns the fuzzer functions for the kubeadm apis.
func Funcs(codecs runtimeserializer.CodecFactory) []interface{} {
	return []interface{}{
		fuzzInitConfiguration,
		fuzzClusterConfiguration,
		fuzzAuditPolicyConfiguration,
		fuzzComponentConfigs,
		fuzzNodeRegistration,
		fuzzLocalEtcd,
		fuzzNetworking,
		fuzzJoinConfiguration,
	}
}

func fuzzInitConfiguration(obj *kubeadm.InitConfiguration, c fuzz.Continue) {
	c.FuzzNoCustom(obj)

	// Pinning values for fields that get defaults if fuzz value is empty string or nil (thus making the round trip test fail)

	// Since ClusterConfiguration never get serialized in the external variant of InitConfiguration,
	// it is necessary to apply external api defaults here to get the round trip internal->external->internal working.
	// More specifically:
	// internal with manually applied defaults -> external object : loosing ClusterConfiguration) -> internal object with automatically applied defaults
	obj.ClusterConfiguration = kubeadm.ClusterConfiguration{
		AuditPolicyConfiguration: kubeadm.AuditPolicyConfiguration{
			LogDir:    constants.StaticPodAuditPolicyLogDir,
			LogMaxAge: &v1beta1.DefaultAuditPolicyLogMaxAge,
		},
		CertificatesDir: v1beta1.DefaultCertificatesDir,
		ClusterName:     v1beta1.DefaultClusterName,
		Etcd: kubeadm.Etcd{
			Local: &kubeadm.LocalEtcd{
				DataDir: v1beta1.DefaultEtcdDataDir,
			},
		},
		ImageRepository:   v1beta1.DefaultImageRepository,
		KubernetesVersion: v1beta1.DefaultKubernetesVersion,
		Networking: kubeadm.Networking{
			ServiceSubnet: v1beta1.DefaultServicesSubnet,
			DNSDomain:     v1beta1.DefaultServiceDNSDomain,
		},
	}
	// Adds the default bootstrap token to get the round working
	obj.BootstrapTokens = []kubeadm.BootstrapToken{
		{
			// Description
			// Expires
			Groups: []string{"foo"},
			// Token
			TTL:    &metav1.Duration{Duration: 1234},
			Usages: []string{"foo"},
		},
	}
}

func fuzzNodeRegistration(obj *kubeadm.NodeRegistrationOptions, c fuzz.Continue) {
	c.FuzzNoCustom(obj)

	// Pinning values for fields that get defaults if fuzz value is empty string or nil (thus making the round trip test fail)
	obj.CRISocket = "foo"
}

func fuzzClusterConfiguration(obj *kubeadm.ClusterConfiguration, c fuzz.Continue) {
	c.FuzzNoCustom(obj)

	// Pinning values for fields that get defaults if fuzz value is empty string or nil (thus making the round trip test fail)
	obj.CertificatesDir = "foo"
	obj.CIImageRepository = "" //This fields doesn't exists in public API >> using default to get the roundtrip test pass
	obj.ClusterName = "bar"
	obj.ImageRepository = "baz"
	obj.KubernetesVersion = "qux"
}

func fuzzAuditPolicyConfiguration(obj *kubeadm.AuditPolicyConfiguration, c fuzz.Continue) {
	c.FuzzNoCustom(obj)

	// Pinning values for fields that get defaults if fuzz value is empty string or nil (thus making the round trip test fail)
	obj.LogDir = "foo"
	obj.LogMaxAge = new(int32)
}

func fuzzComponentConfigs(obj *kubeadm.ComponentConfigs, c fuzz.Continue) {
	// This is intentionally empty because component config does not exists in the public api
	// (empty mean all ComponentConfigs fields nil, and this is necessary for getting roundtrip passing)
}

func fuzzLocalEtcd(obj *kubeadm.LocalEtcd, c fuzz.Continue) {
	c.FuzzNoCustom(obj)

	// Pinning values for fields that get defaults if fuzz value is empty string or nil (thus making the round trip test fail)
	obj.DataDir = "foo"
}

func fuzzNetworking(obj *kubeadm.Networking, c fuzz.Continue) {
	c.FuzzNoCustom(obj)

	// Pinning values for fields that get defaults if fuzz value is empty string or nil (thus making the round trip test fail)
	obj.DNSDomain = "foo"
	obj.ServiceSubnet = "bar"
}

func fuzzJoinConfiguration(obj *kubeadm.JoinConfiguration, c fuzz.Continue) {
	c.FuzzNoCustom(obj)

	// Pinning values for fields that get defaults if fuzz value is empty string or nil (thus making the round trip test fail)
	obj.CACertPath = "foo"
	obj.ClusterName = "bar"
	obj.DiscoveryTimeout = &metav1.Duration{Duration: 1234}
	obj.DiscoveryToken = "baz"
	obj.TLSBootstrapToken = "qux"
}
