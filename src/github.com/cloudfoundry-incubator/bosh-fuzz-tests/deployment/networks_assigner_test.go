package deployment_test

import (
	fakebftdepl "github.com/cloudfoundry-incubator/bosh-fuzz-tests/deployment/fakes"

	. "github.com/cloudfoundry-incubator/bosh-fuzz-tests/deployment"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("NetworksAssigner", func() {
	var (
		networksAssigner NetworksAssigner
		networks         [][]string
	)

	BeforeEach(func() {
		networks = [][]string{[]string{"manual"}}
		nameGenerator := &fakebftdepl.FakeNameGenerator{}
		nameGenerator.Names = []string{"foo-net", "bar-net", "baz-net"}
		networksAssigner = NewSeededNetworksAssigner(networks, nameGenerator, 5)
	})

	It("assigns network of the given type to job and cloud config", func() {
		inputs := []Input{
			{
				Jobs: []Job{
					{
						Name:              "foo",
						Instances:         2,
						AvailabilityZones: []string{"z1"},
					},
				},
				CloudConfig: CloudConfig{
					AvailabilityZones: []string{"z1"},
				},
			},
		}

		networksAssigner.Assign(inputs)

		Expect(inputs).To(Equal([]Input{
			{
				Jobs: []Job{
					{
						Name:              "foo",
						Instances:         2,
						AvailabilityZones: []string{"z1"},
						Networks: []JobNetworkConfig{
							{
								Name:          "foo-net",
								DefaultDNSnGW: true,
							},
						},
					},
				},
				CloudConfig: CloudConfig{
					AvailabilityZones: []string{"z1"},
					Networks: []NetworkConfig{
						{
							Name: "foo-net",
							Type: "manual",
							Subnets: []SubnetConfig{
								{
									IpRange:           "192.168.0.0/24",
									Gateway:           "192.168.0.1",
									AvailabilityZones: []string{"z1"},
									Reserved: []string{
										"192.168.0.9",
										"192.168.0.93",
										"192.168.0.149-192.168.0.159",
									},
								},
							},
						},
						{
							Name: "bar-net",
							Type: "manual",
							Subnets: []SubnetConfig{
								{
									IpRange: "192.168.1.0/24",
									Gateway: "192.168.1.254",
									Reserved: []string{
										"192.168.1.11",
										"192.168.1.120",
										"192.168.1.186-192.168.1.234",
									},
								},
							},
						},
					},
					CompilationNetwork: "bar-net",
				},
			},
		},
		))
	})

	It("generates new subnet range for each subnet", func() {
		inputs := []Input{
			{
				Jobs: []Job{
					{
						Name:              "foo",
						Instances:         1,
						AvailabilityZones: []string{"z1"},
					},
					{
						Name:              "bar",
						Instances:         1,
						AvailabilityZones: []string{"z2"},
					},
				},
				CloudConfig: CloudConfig{
					AvailabilityZones: []string{"z1", "z2"},
				},
			},
		}
		networksAssigner.Assign(inputs)

		Expect(inputs).To(Equal([]Input{
			{
				Jobs: []Job{
					{
						Name:              "foo",
						Instances:         1,
						AvailabilityZones: []string{"z1"},
						Networks: []JobNetworkConfig{
							{
								Name:          "foo-net",
								DefaultDNSnGW: true,
							},
						},
					},
					{
						Name:              "bar",
						Instances:         1,
						AvailabilityZones: []string{"z2"},
						Networks: []JobNetworkConfig{
							{
								Name:          "foo-net",
								DefaultDNSnGW: true,
							},
						},
					},
				},
				CloudConfig: CloudConfig{
					AvailabilityZones: []string{"z1", "z2"},
					Networks: []NetworkConfig{
						{
							Name: "foo-net",
							Type: "manual",
							Subnets: []SubnetConfig{
								{
									IpRange:           "192.168.0.0/24",
									Gateway:           "192.168.0.1",
									AvailabilityZones: []string{"z2"},
									Reserved: []string{
										"192.168.0.45-192.168.0.237",
									},
								},
								{
									IpRange:           "192.168.1.0/24",
									Gateway:           "192.168.1.254",
									AvailabilityZones: []string{"z2", "z1"},
									Reserved: []string{
										"192.168.1.11",
										"192.168.1.120",
										"192.168.1.186-192.168.1.234",
									},
								},
							},
						},
						{
							Name: "bar-net",
							Type: "manual",
							Subnets: []SubnetConfig{
								{
									IpRange:  "192.168.2.0/24",
									Gateway:  "192.168.2.1",
									Reserved: []string{"192.168.2.30"},
								},
								{
									IpRange:  "192.168.3.0/24",
									Gateway:  "192.168.3.254",
									Reserved: []string{"192.168.3.141"},
								},
							},
						},
					},
					CompilationNetwork:          "foo-net",
					CompilationAvailabilityZone: "z1",
				},
			},
		},
		))
	})
})
