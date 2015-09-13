package wundergo_integration_test

import (
	"github.com/nu7hatch/gouuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/robdimsdale/wundergo"
)

const (
	processorType = "generic"
	configuration = ""
)

var _ = Describe("basic webhook functionality", func() {
	var (
		newList wundergo.List
	)

	BeforeEach(func() {
		var err error

		By("Creating a new list")
		uuid1, err := uuid.NewV4()
		Expect(err).NotTo(HaveOccurred())
		newListTitle := uuid1.String()

		newList, err = client.CreateList(newListTitle)
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		var err error

		By("Deleting new list")
		Eventually(func() error {
			newList, err = client.List(newList.ID)
			return client.DeleteList(newList)
		}).Should(Succeed())

		var lists []wundergo.List
		Eventually(func() (bool, error) {
			lists, err = client.Lists()
			return listContains(lists, newList), err
		}).Should(BeFalse())
	})

	It("can list, create and delete webhooks", func() {

		By("Listing existing webhooks")
		webhooks, err := client.WebhooksForListID(newList.ID)
		Expect(err).NotTo(HaveOccurred())
		Expect(len(webhooks)).To(BeZero())

		By("Creating a new webhook")
		url := "https://some-fake-url.com"

		newWebhook, err := client.CreateWebhook(
			newList.ID,
			url,
			processorType,
			configuration,
		)
		Expect(err).NotTo(HaveOccurred())
		Expect(newWebhook.URL).To(Equal(url))

		By("Validating the new webhook is present in list")
		Eventually(func() (bool, error) {
			webhooks, err := client.WebhooksForListID(newList.ID)
			return webhooksContain(webhooks, newWebhook), err
		}).Should(BeTrue())

		By("Deleting the new webhook")
		err = client.DeleteWebhook(newWebhook)
		Expect(err).NotTo(HaveOccurred())

		By("Validating the new webhook is not present in list")
		Eventually(func() (bool, error) {
			webhooks, err := client.WebhooksForListID(newList.ID)
			return webhooksContain(webhooks, newWebhook), err
		}).Should(BeFalse())
	})
})

func webhooksContain(webhooks []wundergo.Webhook, webhook wundergo.Webhook) bool {
	for _, w := range webhooks {
		if w.ID == webhook.ID {
			return true
		}
	}
	return false
}