package integration

import (
	. "github.com/containers/podman/v4/test/utils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
)

var _ = Describe("Podman pod kill", func() {

	It("podman pod kill bogus", func() {
		session := podmanTest.Podman([]string{"pod", "kill", "foobar"})
		session.WaitWithDefaultTimeout()
		Expect(session).To(ExitWithError())
	})

	It("podman pod kill a pod by id", func() {
		_, ec, podid := podmanTest.CreatePod(nil)
		Expect(ec).To(Equal(0))

		session := podmanTest.RunTopContainerInPod("", podid)
		session.WaitWithDefaultTimeout()
		Expect(session).Should(ExitCleanly())

		session = podmanTest.RunTopContainerInPod("", podid)
		session.WaitWithDefaultTimeout()
		Expect(session).Should(ExitCleanly())

		result := podmanTest.Podman([]string{"pod", "kill", podid})
		result.WaitWithDefaultTimeout()
		Expect(result).Should(ExitCleanly())
		Expect(podmanTest.NumberOfContainersRunning()).To(Equal(0))
	})

	It("podman pod kill a pod by id with TERM", func() {
		_, ec, podid := podmanTest.CreatePod(nil)
		Expect(ec).To(Equal(0))

		session := podmanTest.RunTopContainerInPod("", podid)
		session.WaitWithDefaultTimeout()
		Expect(session).Should(ExitCleanly())

		result := podmanTest.Podman([]string{"pod", "kill", "-s", "9", podid})
		result.WaitWithDefaultTimeout()
		Expect(result).Should(ExitCleanly())
		Expect(podmanTest.NumberOfContainersRunning()).To(Equal(0))
	})

	It("podman pod kill a pod by name", func() {
		_, ec, podid := podmanTest.CreatePod(map[string][]string{"--name": {"test1"}})
		Expect(ec).To(Equal(0))

		session := podmanTest.RunTopContainerInPod("", podid)
		session.WaitWithDefaultTimeout()
		Expect(session).Should(ExitCleanly())

		result := podmanTest.Podman([]string{"pod", "kill", "test1"})
		result.WaitWithDefaultTimeout()
		Expect(result).Should(ExitCleanly())
		Expect(podmanTest.NumberOfContainersRunning()).To(Equal(0))
	})

	It("podman pod kill a pod by id with a bogus signal", func() {
		_, ec, podid := podmanTest.CreatePod(map[string][]string{"--name": {"test1"}})
		Expect(ec).To(Equal(0))

		session := podmanTest.RunTopContainerInPod("", podid)
		session.WaitWithDefaultTimeout()
		Expect(session).Should(ExitCleanly())

		result := podmanTest.Podman([]string{"pod", "kill", "-s", "bogus", "test1"})
		result.WaitWithDefaultTimeout()
		Expect(result).Should(Exit(125))
		Expect(podmanTest.NumberOfContainersRunning()).To(Equal(1))
	})

	It("podman pod kill latest pod", func() {
		_, ec, podid := podmanTest.CreatePod(nil)
		Expect(ec).To(Equal(0))

		session := podmanTest.RunTopContainerInPod("", podid)
		session.WaitWithDefaultTimeout()
		Expect(session).Should(ExitCleanly())

		_, ec, podid2 := podmanTest.CreatePod(nil)
		Expect(ec).To(Equal(0))

		session = podmanTest.RunTopContainerInPod("", podid2)
		session.WaitWithDefaultTimeout()
		Expect(session).Should(ExitCleanly())

		session = podmanTest.RunTopContainerInPod("", podid2)
		session.WaitWithDefaultTimeout()
		Expect(session).Should(ExitCleanly())
		if !IsRemote() {
			podid2 = "-l"
		}
		result := podmanTest.Podman([]string{"pod", "kill", podid2})
		result.WaitWithDefaultTimeout()
		Expect(result).Should(ExitCleanly())
		Expect(podmanTest.NumberOfContainersRunning()).To(Equal(1))
	})

	It("podman pod kill all", func() {
		SkipIfRootlessCgroupsV1("Not supported for rootless + CgroupsV1")
		_, ec, podid := podmanTest.CreatePod(nil)
		Expect(ec).To(Equal(0))

		session := podmanTest.RunTopContainerInPod("", podid)
		session.WaitWithDefaultTimeout()
		Expect(session).Should(ExitCleanly())

		session = podmanTest.RunTopContainerInPod("", podid)
		session.WaitWithDefaultTimeout()
		Expect(session).Should(ExitCleanly())

		_, ec, podid2 := podmanTest.CreatePod(nil)
		Expect(ec).To(Equal(0))

		session = podmanTest.RunTopContainerInPod("", podid2)
		session.WaitWithDefaultTimeout()
		Expect(session).Should(ExitCleanly())

		result := podmanTest.Podman([]string{"pod", "kill", "-a"})
		result.WaitWithDefaultTimeout()
		GinkgoWriter.Println(result.OutputToString(), result.ErrorToString())
		Expect(result).Should(ExitCleanly())
		Expect(podmanTest.NumberOfContainersRunning()).To(Equal(0))
	})
})
