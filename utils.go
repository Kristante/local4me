package main

import (
	"fmt"
	"sort"
	"time"
)

var Names = map[string]string{
	"Иванов Дмитрий Николаевич":         "@kristante",
	"Булавина Василина Васильевна":      "@vslnb",
	"Дюжов Артём Витальевич":            "Заявка Артёма",
	"Епифановский Михаил Александрович": "Заявка Миши",
	"Ершов Александр Павлович":          "Заявка Саши",
	"Любимов Георгий Владимирович":      "Заявка Гоши",
	"Потриваев Никита Андреевич":        "Заявка Никиты",
	"Хасаншин Рустам Альбертович":       "Заявка Рустама",
	"Хусниярова Алия Раисовна":          "Заявка Алии",
	"Cлужебная УЗ":                      "",
	"Автоматизация":                     "",
}

func CheckMemberName(request Requests) bool {
	if request.Member.Name == "" || request.Team.Name == "Техподдержка" {
		return true
	}
	return false
}

func CheckExcludedNames(name string) bool {
	excludedNames := Names
	if _, exists := excludedNames[name]; exists {
		return false
	}
	return true
}

func ConvertInfoForMessageTelegram(request Request) string {
	return fmt.Sprintf("Появилась новая заявка под номером %d от пользователя %s\nОзнакомиться подробнее можно по ссылке: https://rpa.itsm.mos.ru/requests/%d", request.ID, request.CreatedBy.Name, request.ID)
}

func ConvertNotesForMessageTelegram(note Note, reqID int) string {
	link := note.Person.Name
	return fmt.Sprintf("%s\nПоявился новый комментарий от пользователя %s у заявки под номером %d:\nТекст комментария: %s\nОзнакомиться подробнее: https://rpa.itsm.mos.ru/requests/%d", link, note.Person.Name, reqID, note.Text, reqID)
}

func GetComments(Notes []Note) *Note {
	// Необходимо отсортировать по дате создания по убыванию
	sort.Slice(Notes, func(i, j int) bool {
		t1, _ := time.Parse(time.RFC3339, Notes[i].CreatedAt)
		t2, _ := time.Parse(time.RFC3339, Notes[j].CreatedAt)
		return t1.After(t2)
	})
	if len(Notes) > 0 {
		latestNote := Notes[0]
		if CheckExcludedNames(latestNote.Person.Name) {
			return &latestNote
		}
		fmt.Println("Имя создателя входит в список исключений, выводить информацию не нужно.")
	}
	return nil
}
