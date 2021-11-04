package harvester

import (
	"context"
	"fmt"

	"github.com/tinkerbell/boots/ipxe"
	"github.com/tinkerbell/boots/job"
)

type Installer struct{}

func (i Installer) BootScriptHarvester020() job.BootScript {
	return func(ctx context.Context, j job.Job, s ipxe.Script) ipxe.Script {
		return bootScriptHarvester(ctx, j, s, "v0.2.0")
	}
}

func (i Installer) BootScriptHarvester030() job.BootScript {
	return func(ctx context.Context, j job.Job, s ipxe.Script) ipxe.Script {
		return bootScriptHarvester(ctx, j, s, "v0.3.0")
	}
}
func bootScriptHarvester(ctx context.Context, j job.Job, s ipxe.Script, version string) ipxe.Script {
	s.PhoneHome("provisioning.104.01")
	if len(j.OSIEBaseURL()) != 0 {
		s.Set("base-url", j.OSIEBaseURL())
	} else {
		s.Set("base-url", "https://releases.rancher.com/harvester")
	}
	j.With("parsed userdata", j.UserData())
	ks := kernelParams(j, s, version)
	if version == "v0.2.0" {
		ks.Initrd(fmt.Sprintf("${base-url}/%s/harvester-initrd-amd64", version))
	} else {
		ks.Initrd(fmt.Sprintf("${base-url}/%s/harvester-%s-initrd-amd64", version, version))
	}
	ks.Boot()
	// once boot script is served no long want this, since harvester install triggers a
	// reboot of the node
	j.DisablePXE(ctx)

	return ks
}

func kernelParams(j job.Job, s ipxe.Script, version string) ipxe.Script {
	currentUserData := j.UserData()

	// we will check userdata to see if flags exist else we add the info //

	switch version {
	case "v0.2.0":
		s.Kernel(fmt.Sprintf("${base-url}/%s/harvester-vmlinuz-amd64", version))
		s.Args("k3os.mode=install", "k3os.debug", "console=tty1,115200", "harvester.install.automatic=true", "boot_cmd=\"echo include_ping_test=yes >> /etc/conf.d/net-online\"")
	case "v0.3.0":
		s.Kernel(fmt.Sprintf("${base-url}/%s/harvester-%s-vmlinuz-amd64", version, version))
		s.Args("rd.cos.disable", "rd.noverifyssl", "net.ifnames=1", "console=tty1", "harvester.install.automatic=true", "boot_cmd=\"echo include_ping_test=yes >> /etc/conf.d/net-online\"")
		s.Args(fmt.Sprintf("root=live:${base-url}/%s/harvester-%s-rootfs-amd64.squashfs", version, version))
	}
	if len(currentUserData) != 0 {
		s.Args(currentUserData)
	}

	return s
}
