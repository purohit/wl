package wundergo_integration_test

import (
	"github.com/nu7hatch/gouuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/robdimsdale/wundergo"
)

var _ = Describe("basic reminder functionality", func() {
	var (
		newList wundergo.List
		newTask wundergo.Task
		err     error
	)

	BeforeEach(func() {
		By("Creating a new list")
		uuid1, err := uuid.NewV4()
		Expect(err).NotTo(HaveOccurred())
		newListTitle := uuid1.String()

		newList, err = client.CreateList(newListTitle)
		Expect(err).NotTo(HaveOccurred())

		By("Creating task in new list")
		uuid, err := uuid.NewV4()
		Expect(err).NotTo(HaveOccurred())
		newTaskTitle := uuid.String()

		Eventually(func() error {
			newTask, err = client.CreateTask(
				newTaskTitle,
				newList.ID,
				0,
				false,
				"",
				0,
				"1970-01-01",
				false,
			)
			return err
		}).ShouldNot(HaveOccurred())
	})

	AfterEach(func() {
		By("Deleting task")
		Eventually(func() error {
			newTask, err = client.Task(newTask.ID)
			return client.DeleteTask(newTask)
		}).Should(Succeed())

		var tasks []wundergo.Task
		Eventually(func() (bool, error) {
			tasks, err = client.TasksForListID(newList.ID)
			return taskContains(tasks, newTask), err
		}).Should(BeFalse())

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

	It("can perform reminder CRUD", func() {
		By("Creating reminder")
		var reminder wundergo.Reminder
		reminderDate := "1970-08-30T08:29:46.203Z"
		createdByDeviceUdid := ""
		Eventually(func() error {
			reminder, err = client.CreateReminder(reminderDate, newTask.ID, createdByDeviceUdid)
			return err
		}).Should(Succeed())

		By("Verifying reminder exists in all reminders")
		Eventually(func() (bool, error) {
			allReminders, err := client.Reminders()
			return reminderContains(allReminders, reminder), err
		}).Should(BeTrue())

		By("Verifying reminder exists in reminders for list")
		Eventually(func() (bool, error) {
			remindersForList, err := client.RemindersForListID(newList.ID)
			return reminderContains(remindersForList, reminder), err
		}).Should(BeTrue())

		By("Verifying reminder exists in reminders for task")
		Eventually(func() (bool, error) {
			remindersForTask, err := client.RemindersForTaskID(newTask.ID)
			return reminderContains(remindersForTask, reminder), err
		}).Should(BeTrue())

		By("Updating reminder")
		reminder.Date = "1971-08-30T08:29:46.203Z"
		Eventually(func() error {
			reminder, err = client.UpdateReminder(reminder)
			return err
		}).Should(Succeed())

		By("Getting reminder")
		var aReminder wundergo.Reminder
		Eventually(func() error {
			aReminder, err = client.Reminder(reminder.ID)
			return err
		}).Should(Succeed())

		Expect(aReminder.ID).To(Equal(reminder.ID))
		Expect(aReminder.Date).To(Equal(reminder.Date))

		By("Deleting reminder")
		Eventually(func() error {
			return client.DeleteReminder(reminder)
		}).Should(Succeed())
	})
})

func reminderContains(reminders []wundergo.Reminder, reminder wundergo.Reminder) bool {
	for _, n := range reminders {
		if n.ID == reminder.ID {
			return true
		}
	}
	return false
}
