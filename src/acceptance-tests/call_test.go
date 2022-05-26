package acceptance_test

import (
	"regexp"
	"strings"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/cloudfoundry/cf-acceptance-tests/helpers/random_name"
	"github.com/cloudfoundry/cf-test-helpers/cf"
)

var _ = Describe("Scheduler Calls", func() {
	var (
		callName string
	)

	BeforeEach(func() {
		callName = random_name.CATSRandomName("CALL")
	})

	AfterEach(func() {
		Expect(cf.Cf("delete-call", callName).Wait(time.Second * 10).Out.Contents()).
			Should(ContainSubstring("Deleted call")) // FIX ME
	})

	Describe("create-call", func() {
		It("test correct call creation", func() {
			Expect(cf.Cf("create-call", appName, callName, `https://www.starkandwayne.com/`).
				Wait(time.Second * 10).Out.Contents()).
				Should(ContainSubstring("OK"))
		})
	})

	Describe("schedule-call", func() {
		It("test correct call scheduling", func() {
			Expect(cf.Cf("create-call", appName, callName, `https://www.starkandwayne.com/`).
				Wait(time.Second * 10).Out.Contents()).
				Should(ContainSubstring("OK"))

			Expect(cf.Cf("schedule-call", callName, `15 * * * *`).
				Wait(time.Second * 10).Out.Contents()).
				Should(ContainSubstring(callName))

			Expect(cf.Cf("call-schedules").
				Wait(time.Second * 10).Out.Contents()).
				Should(ContainSubstring(callName))
		})
	})

	Describe("run-call", func() {
		It("test correct call manual execution", func() {
			Expect(cf.Cf("create-call", appName, callName, `https://www.starkandwayne.com/`).
				Wait(time.Second * 10).Out.Contents()).
				Should(ContainSubstring("OK"))

			Expect(cf.Cf("run-call", callName).
				Wait(time.Second * 10).Out.Contents()).
				Should(ContainSubstring("OK"))
		})
	})

	Describe("delete-call", func() {
		It("test correct call deletion", func() {
			Expect(cf.Cf("create-call", appName, callName, `https://www.starkandwayne.com/`).
				Wait(time.Second * 10).Out.Contents()).
				Should(ContainSubstring("OK"))

			callName := random_name.CATSRandomName("CALL")
			Expect(cf.Cf("create-call", appName, callName, `https://www.starkandwayne.com/`).
				Wait(time.Second * 10).Out.Contents()).
				Should(ContainSubstring("OK"))

			Expect(cf.Cf("delete-call", callName).
				Wait(time.Second * 10).Out.Contents()).
				Should(ContainSubstring("Deleted call")) // FIX ME

			Expect(cf.Cf("calls").
				Wait(time.Second * 10).Out.Contents()).
				ShouldNot(ContainSubstring(callName))

			// Expect(cf.Cf("call-history", callName).
			// 	Wait(time.Second * 10).Out.Contents()).
			// 	ShouldNot(ContainSubstring(callName))
		})
	})

	Describe("delete-call-schedule", func() {
		It("test correct call schedule deletion", func() {
			Expect(cf.Cf("create-call", appName, callName, `https://www.starkandwayne.com/`).
				Wait(time.Second * 10).Out.Contents()).
				Should(ContainSubstring("OK"))

			Expect(cf.Cf("schedule-call", callName, `15 * * * *`).
				Wait(time.Second * 10).Out.Contents()).
				Should(ContainSubstring(callName))

			schedules := cf.Cf("call-schedules").
				Wait(time.Second * 10)

			Expect(schedules.Out.Contents()).
				Should(ContainSubstring(callName))

			var schedule string
			re := regexp.MustCompile(`^(.*?)[\s]+(.*?)[\s]+([a-z0-9]{8}-[a-z0-9]{4}-[a-z0-9]{4}-[a-z0-9]{4}-[a-z0-9]{12})`)
			for _, line := range strings.Split(string(schedules.Out.Contents()), "\n") {
				for _, i := range re.FindAllStringSubmatch(line, -1) {
					if i[1] == callName {
						schedule = i[3]
					}
				}
			}

			Expect(schedule).NotTo(BeEmpty())

			Expect(cf.Cf("delete-call-schedule", callName, schedule).
				Wait(time.Second * 10).Out.Contents()).
				Should(ContainSubstring(schedule + " deleted")) // FIX ME

			Expect(cf.Cf("call-schedules").
				Wait(time.Second * 10).Out.Contents()).
				ShouldNot(ContainSubstring(schedule))
		})
	})
})
