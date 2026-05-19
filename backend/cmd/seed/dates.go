package main

import "time"

// Фиксированные даты для воспроизводимых результатов сидинга. Все значения
// в UTC. Сидер не выставляет created_at/updated_at в now(), чтобы данные
// не зависели от текущего системного времени.

func ts(month time.Month, day, hour, minute int) time.Time {
	return time.Date(2026, month, day, hour, minute, 0, 0, time.UTC)
}

// Стандартные ГОСТ-шаблоны библиотеки. created_at разнесены по 2026 году,
// у части шаблонов есть updated_at.
var (
	dateGOSTASCreated = ts(time.January, 22, 11, 35)
	dateGOSTASUpdated = mustPtr(ts(time.April, 9, 16, 20))

	dateGOSTProgramCreated = ts(time.February, 18, 14, 8)

	dateGOSTPMICreated = ts(time.March, 26, 10, 42)
	dateGOSTPMIUpdated = mustPtr(ts(time.May, 7, 13, 55))

	dateGOSTOperatorCreated = ts(time.May, 12, 16, 15)
)

// Кастомные шаблоны и связанные с ними версии.
var (
	datePassportCreated = ts(time.April, 8, 11, 17)
	datePassportVersion = ts(time.April, 8, 11, 24)
	datePassportUpdated = mustPtr(ts(time.May, 6, 18, 42))

	dateReleaseCardCreated = ts(time.April, 22, 15, 33)
	dateReleaseCardVersion = ts(time.April, 22, 15, 47)

	dateADRCreated = ts(time.May, 12, 9, 54)
	dateADRVersion = ts(time.May, 12, 10, 5)
)

// Задачи под кастомные шаблоны.
var (
	dateTaskSuccessCreated = ts(time.May, 15, 10, 24)
	dateTaskSuccessUpdated = ts(time.May, 15, 10, 24).Add(7 * time.Second)

	dateTaskFailCreated = ts(time.May, 17, 16, 51)
	dateTaskFailUpdated = ts(time.May, 17, 16, 51).Add(5 * time.Second)
)

// Шаблоны полноформатной документации по NotifyHub и задачи к ним.
var (
	dateApprobationTZCreated  = ts(time.April, 15, 9, 30)
	dateApprobationTZVersion  = ts(time.April, 15, 14, 12)
	dateApprobationTZTask     = ts(time.April, 16, 11, 5)
	dateApprobationTZTaskDone = ts(time.April, 16, 11, 5).Add(11 * time.Second)

	dateApprobationPMICreated  = ts(time.April, 25, 13, 8)
	dateApprobationPMIVersion  = ts(time.April, 25, 18, 41)
	dateApprobationPMITask     = ts(time.May, 2, 10, 17)
	dateApprobationPMITaskDone = ts(time.May, 2, 10, 17).Add(8 * time.Second)

	dateApprobationOPCreated  = ts(time.May, 8, 12, 22)
	dateApprobationOPVersion  = ts(time.May, 8, 16, 3)
	dateApprobationOPTask     = ts(time.May, 10, 9, 41)
	dateApprobationOPTaskDone = ts(time.May, 10, 9, 41).Add(9 * time.Second)
)

func mustPtr(t time.Time) *time.Time { return &t }
