package main

import "github.com/vembry/saga"

func main() {

	saga := saga.New()

	saga.Start()
	saga.Stop()

	// activity := saga.NewActivity("activity0", commit, rollback)
	// activity1 := saga.NewActivity("activity1", commit, rollback)
	// activity2 := saga.NewActivity("activity2", commit, rollback)
	// activity3 := saga.NewActivity("activity3", commit, rollback)
	// activity4 := saga.NewActivity("activity4", commit, rollback)
	// activity5 := saga.NewActivity("activity5", commit, rollback)

	// workflow := saga.NewWorkflow[mockparameter](
	// 	"workflow0",
	// 	activity,
	// 	activity1,
	// 	activity2,
	// 	activity3,
	// 	activity4,
	// 	activity5,
	// )

	// workflow.Execute(mockparameter{})
}

// type mockparameter struct {
// 	Value []string
// }

// func commit(arg mockparameter) {
// 	log.Printf("committing...")
// 	arg.Value = append(arg.Value, fmt.Sprintf("%d", time.Now().UnixMilli()))
// }

// func rollback(arg mockparameter) {
// 	log.Printf("rolling back...")
// 	arg.Value = arg.Value[0 : len(arg.Value)-1]
// }
