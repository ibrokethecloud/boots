package ros2

import (
	"context"
	"github.com/tinkerbell/boots/ipxe"
	"github.com/tinkerbell/boots/job"
)

type Installer struct{}

func (i Installer) BootScript() job.BootScript {
	return func(ctx context.Context, j job.Job, s ipxe.Script) ipxe.Script {
		s.PhoneHome("provisioning.104.01")
		s.Set("arch", j.Arch())
		s.Set("version", j.OperatingSystem().Version)
		s.Set("url", "https://github.com/rancher-sandbox/os2/releases/download/${version}")
		s.Set("kernel", "rancheros-${version}-${arch}-kernel")
		s.Set("initrd", "rancheros-${version}-${arch}-initrd")
		s.Set("rootfs", "rancheros-${version}-${arch}.squashfs")
		s.Set("iso", "rancheros-${version}-${arch}.iso")
		s.Kernel("${base-url}/vmlinuz")
		s.Args("${url}/${kernel}", "initrd=${initrd}", "ip=dhcp", "rd.cos.disable", "root=live:${url}/${rootfs}",
			"rancheros.install.automatic=true", "rancheros.install.iso_url=${url}/${iso}",
			"rancheros.install.config_url=${config}", "console=tty1")
		s.Args(j.UserData())
		s.Initrd("rancheros-${version}-${arch}-initrd")
		s.Boot()
		j.DisablePXE(ctx)
		return s
	}
}

/* Sample pxe boot script
set arch amd64
set version v0.1.0-alpha16
set url https://github.com/rancher-sandbox/os2/releases/download/${version}
set kernel rancheros-${version}-${arch}-kernel
set initrd rancheros-${version}-${arch}-initrd
set rootfs rancheros-${version}-${arch}.squashfs
set iso    rancheros-${version}-${arch}.iso
# set config http://example.com/machine-config
# set cmdline extra.values=1
kernel ${url}/${kernel} initrd=${initrd} ip=dhcp rd.cos.disable root=live:${url}/${rootfs} rancheros.install.automatic=true rancheros.install.iso_url=${url}/${iso} rancheros.install.config_url=${config} console=tty1 console=ttyS0 ${cmdline}
initrd ${url}/${initrd}
boot
*/
