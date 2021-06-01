package harvester

import (
	"github.com/tinkerbell/boots/ipxe"
	"github.com/tinkerbell/boots/job"
)

func init() {
	job.RegisterSlug("harvester_0_2_0", bootScriptHarvester)

}

func bootScriptHarvester(j job.Job, s *ipxe.Script) {
	s.PhoneHome("provisioning.104.01")
	if len(j.OSIEBaseURL()) != 0 {
		s.Set("base-url", j.OSIEBaseURL())
	} else {
		s.Set("base-url", "https://releases.rancher.com/harvester/master")
	}

	s.Kernel("${base-url}/harvester-vmlinuz-amd64")

	j.With("parsed userdata", j.UserData())
	kernelParams(j, s)

	s.Initrd("${base-url}/harvester-initrd-amd64")
	s.Boot()
	// once boot script is served no long want this, since harvester install triggers a
	// reboot of the node
	j.DisablePXE()
}

func kernelParams(j job.Job, s *ipxe.Script) {
	s.Args("k3os.mode=install", "k3os.debug", "console=tty1,115200", "harvester.install.automatic=true", "boot_cmd=\"echo include_ping_test=yes >> /etc/conf.d/net-online\"")
	if len(j.UserData()) != 0 {
		s.Args(j.UserData())
	}
}
