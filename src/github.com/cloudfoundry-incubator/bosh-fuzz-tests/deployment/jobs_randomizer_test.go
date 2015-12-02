package deployment_test

import (
	"math/rand"

	bftconfig "github.com/cloudfoundry-incubator/bosh-fuzz-tests/config"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"

	. "github.com/cloudfoundry-incubator/bosh-fuzz-tests/deployment"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("JobsRandomizer", func() {
	var (
		jobsRandomizer JobsRandomizer
	)

	It("generates extra input for migrated jobs", func() {
		parameters := bftconfig.Parameters{
			NameLength:               []int{5},
			Instances:                []int{2},
			AvailabilityZones:        [][]string{[]string{"z1"}, []string{"z1", "z2"}},
			PersistentDiskDefinition: []string{"disk_pool"},
			PersistentDiskSize:       []int{100},
			NumberOfJobs:             []int{2},
			MigratedFromCount:        []int{0},
			VmTypeDefinition:         []string{"vm_type"},
			StemcellDefinition:       []string{"os_version"},
		}
		logger := boshlog.NewLogger(boshlog.LevelNone)
		rand.Seed(64)
		nameGenerator := NewNameGenerator()
		jobsRandomizer = NewJobsRandomizer(parameters, 2, nameGenerator, logger)

		inputs, err := jobsRandomizer.Generate()
		Expect(err).ToNot(HaveOccurred())

		Expect(inputs).To(Equal([]Input{
			{
				Jobs: []Job{
					{
						Name:               "joNAw",
						Instances:          2,
						AvailabilityZones:  []string{"z1"},
						PersistentDiskPool: "jO2GYdmz6a",
						VmType:             "yRjDTBs3VX",
					},
					{
						Name:               "gQ8el",
						Instances:          2,
						AvailabilityZones:  []string{"z1"},
						PersistentDiskPool: "b0xNg3RWDs",
						VmType:             "mO5Qu91qDq",
					},
				},
				CloudConfig: CloudConfig{
					AvailabilityZones: []string{"z1"},
					PersistentDiskPools: []DiskConfig{
						{Name: "jO2GYdmz6a", Size: 100},
						{Name: "b0xNg3RWDs", Size: 100},
					},
					VmTypes: []VmTypeConfig{
						{Name: "yRjDTBs3VX"},
						{Name: "mO5Qu91qDq"},
					},
				},
				Stemcells: []StemcellConfig{
					{OS: "toronto-os", Version: "1", Alias: "default"},
				},
			},
			{
				Jobs: []Job{
					{
						Name:               "joNAw",
						Instances:          2,
						AvailabilityZones:  []string{"z1"},
						PersistentDiskPool: "aREws5gmsY",
						VmType:             "eE7YhmI4yV",
					},
					{
						Name:               "gQ8el",
						Instances:          2,
						AvailabilityZones:  []string{"z1"},
						PersistentDiskPool: "grcWVDVTZN",
						VmType:             "sOgP3i7apW",
					},
				},
				CloudConfig: CloudConfig{
					AvailabilityZones: []string{"z1"},
					PersistentDiskPools: []DiskConfig{
						{Name: "aREws5gmsY", Size: 100},
						{Name: "grcWVDVTZN", Size: 100},
					},
					VmTypes: []VmTypeConfig{
						{Name: "eE7YhmI4yV"},
						{Name: "sOgP3i7apW"},
					},
				},
				Stemcells: []StemcellConfig{
					{OS: "toronto-os", Version: "1", Alias: "default"},
				},
			},
		}))
	})

	It("when migrated job does not have az it sets random az in migrated_from", func() {
		parameters := bftconfig.Parameters{
			NameLength:               []int{5},
			Instances:                []int{2},
			AvailabilityZones:        [][]string{[]string{"z1"}, nil},
			PersistentDiskDefinition: []string{"persistent_disk_size"},
			PersistentDiskSize:       []int{0},
			NumberOfJobs:             []int{1},
			MigratedFromCount:        []int{1},
			VmTypeDefinition:         []string{"vm_type"},
			StemcellDefinition:       []string{"name"},
		}
		logger := boshlog.NewLogger(boshlog.LevelNone)
		rand.Seed(64)
		nameGenerator := NewNameGenerator()
		jobsRandomizer = NewJobsRandomizer(parameters, 1, nameGenerator, logger)

		inputs, err := jobsRandomizer.Generate()
		Expect(err).ToNot(HaveOccurred())

		Expect(inputs).To(Equal([]Input{
			{
				Jobs: []Job{
					{
						Name:      "vmz6agRjDT",
						Instances: 2,
						VmType:    "rYND0xNg3R",
					},
				},
				CloudConfig: CloudConfig{
					VmTypes: []VmTypeConfig{
						{Name: "rYND0xNg3R"},
					},
				},
				Stemcells: []StemcellConfig{
					{
						Name:    "ubuntu-stemcell",
						Version: "1",
						Alias:   "default",
					},
				},
			},
			{
				Jobs: []Job{
					{
						Name:              "joNAw",
						Instances:         2,
						AvailabilityZones: []string{"z1"},
						VmType:            "arKicN3O2G",
						MigratedFrom: []MigratedFromConfig{
							{Name: "vmz6agRjDT", AvailabilityZone: "z1"},
						},
					},
				},
				CloudConfig: CloudConfig{
					AvailabilityZones: []string{"z1"},
					VmTypes: []VmTypeConfig{
						{Name: "arKicN3O2G"},
					},
				},
				Stemcells: []StemcellConfig{
					{
						Name:    "ubuntu-stemcell",
						Version: "1",
						Alias:   "default",
					},
				},
			},
		}))
	})
})
