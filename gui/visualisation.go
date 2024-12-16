package gui

import (
	e "curse_work/encryption"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"strconv"
	"strings"
)

var currentStep [4]int

var helpers [][]string
var tables [16][8]string
var mainKey string
var mainText string
var content *fyne.Container

//var spaceLabel *wigdet.Label

func Initialization(w fyne.Window) {
	inputText := widget.NewEntry()
	inputText.SetPlaceHolder("Открытый текст (значение имеют только первые 8 символов)")

	keyEntry := widget.NewEntry()
	keyEntry.SetPlaceHolder("Ключ шифрования (ровно 8 символов)")

	startButton := widget.NewButton("Начать шифрование", func() {
		mainText = inputText.Text
		mainKey = keyEntry.Text

		if len(mainText) == 0 {
			dialog.ShowError(fmt.Errorf("введите открытый текст"), w)
			return
		}

		if len(mainKey) < 8 {
			dialog.ShowError(fmt.Errorf("ключ должен быть длиной не менее 8 символов"), w)
			return
		}
		Visualisation(w)
	})

	exitButton := widget.NewButton("Вернуться на главный экран", func() {
		MainMenu(w)
	})

	content := container.NewVBox(
		widget.NewLabel("DES Step-by-Step Visualization"),
		inputText,
		keyEntry,
		startButton,
		exitButton,
	)

	w.SetContent(content)
	w.Resize(fyne.NewSize(800, 600))

}

func Visualisation(w fyne.Window) {
	helpers = [][]string{
		{"Разделение на левую и правую части",
			"Расширение правой части с 32 бит до 48 согласно таблице ET",
			"Исключающее ИЛИ с ключем",
			"S-подстановка",
			"Перестановка левой и правой частей"},
		{"S-подстановка. Выбранное 4-битное значение из таблицы определяется значением подменяемого 6-битного блока", "Перестановка по таблице PT"},
		{"Расширение ключа по таблице PC1",
			"Левый сдвиг левой и правой части по таблице LeftShift",
			"Сжатие ключа по таблице PC2"},
		{"Начальная перестановка текста по таблице IP", "Добавление заполнения до 64 бит", "Сеть Фейстеля",
			"Конечная перестановка текста по таблице FP", "Генерация ключей"},
		{"", "Таблица IP:", "Таблица PT:", "Таблица PC1:",
			"Таблица LeftShift:", "Таблица PC2:", "Таблица FP:"},
		{"Указатель:"}}
	tables = [16][8]string{
		{"", "", "", "", "", "", "", ""}, //null
		{
			"58, 50, 42, 34, 26, 18, 10, 2", //IP
			"60, 52, 44, 36, 28, 20, 12, 4",
			"62, 54, 46, 38, 30, 22, 14, 6",
			"64, 56, 48, 40, 32, 24, 16, 8",
			"57, 49, 41, 33, 25, 17, 9, 1",
			"59, 51, 43, 35, 27, 19, 11, 3",
			"61, 53, 45, 37, 29, 21, 13, 5",
			"63, 55, 47, 39, 31, 23, 15, 7",
		},
		{
			"32, 1, 2, 3, 4, 5", //ET
			"4, 5, 6, 7, 8, 9",
			"8, 9, 10, 11, 12, 13",
			"12, 13, 14, 15, 16, 17",
			"16, 17, 18, 19, 20, 21",
			"20, 21, 22, 23, 24, 25",
			"24, 25, 26, 27, 28, 29",
			"28, 29, 30, 31, 32, 1"},
		{
			"16, 7, 20, 21, 29, 12, 28, 17", //PT
			"1, 15, 23, 26, 5, 18, 31, 10",
			"2, 8, 24, 14, 32, 27, 3, 9",
			"19, 13, 30, 6, 22, 11, 4, 25", "", "", "", ""},
		{
			"57, 49, 41, 33, 25, 17, 9", //PC1
			"1, 58, 50, 42, 34, 26, 18",
			"10, 2, 59, 51, 43, 35, 27",
			"19, 11, 3, 60, 52, 44, 36",
			"63, 55, 47, 39, 31, 23, 15",
			"7, 62, 54, 46, 38, 30, 22",
			"14, 6, 61, 53, 45, 37, 29",
			"21, 13, 5, 28, 20, 12, 4",
		},
		{
			"14, 17, 11, 24, 1, 5, 3, 28", //PC2
			"15, 6, 21, 10, 23, 19, 12, 4",
			"26, 8, 16, 7, 27, 20, 13, 2",
			"41, 52, 31, 37, 47, 55, 30, 40",
			"51, 45, 33, 48, 44, 49, 39, 56",
			"34, 53, 46, 42, 50, 36, 29, 32", ""},
		{
			"40, 8, 48, 16, 56, 24, 64, 32", //FP
			"39, 7, 47, 15, 55, 23, 63, 31",
			"38, 6, 46, 14, 54, 22, 62, 30",
			"37, 5, 45, 13, 53, 21, 61, 29",
			"36, 4, 44, 12, 52, 20, 60, 28",
			"35, 3, 43, 11, 51, 19, 59, 27",
			"34, 2, 42, 10, 50, 18, 58, 26",
			"33, 1, 41, 9, 49, 17, 57, 25",
		},
		{"1, 1, 2, 2, 2, 2, 2, 2, 1, 2, 2, 2, 2, 2, 2, 1", "", "",
			"", "", "", "", ""}, //LeftShifts
		// S-Box 1
		{
			"14, 4, 13, 1, 2, 15, 11, 8, 3, 10, 6, 12, 5, 9, 0, 7",
			"0, 15, 7, 4, 14, 2, 13, 1, 10, 6, 12, 11, 9, 5, 3, 8",
			"4, 1, 14, 8, 13, 6, 2, 11, 15, 12, 9, 7, 3, 10, 5, 0",
			"15, 12, 8, 2, 4, 9, 1, 7, 5, 11, 3, 14, 10, 0, 6, 13", "", "", "", ""},
		// S-Box 2
		{
			"15, 1, 8, 14, 6, 11, 3, 4, 9, 7, 2, 13, 12, 0, 5, 10",
			"3, 13, 4, 7, 15, 2, 8, 14, 12, 0, 1, 10, 6, 9, 11, 5",
			"0, 14, 7, 11, 10, 4, 13, 1, 5, 8, 12, 6, 9, 3, 2, 15",
			"13, 8, 10, 1, 3, 15, 4, 2, 11, 6, 7, 12, 0, 5, 14, 9", "", "", "", ""},
		// S-Box 3
		{
			"10, 0, 9, 14, 6, 3, 15, 5, 1, 13, 12, 7, 11, 4, 2, 8",
			"13, 7, 0, 9, 3, 4, 6, 10, 2, 8, 5, 14, 12, 11, 15, 1",
			"13, 6, 4, 9, 8, 15, 3, 0, 11, 1, 2, 12, 5, 10, 14, 7",
			"1, 10, 13, 0, 6, 9, 8, 7, 4, 15, 14, 3, 11, 5, 2, 12", "", "", "", ""},
		// S-Box 4
		{
			"7, 13, 14, 3, 0, 6, 9, 10, 1, 2, 8, 5, 11, 12, 4, 15",
			"13, 8, 11, 5, 6, 15, 0, 3, 4, 7, 2, 12, 1, 10, 14, 9",
			"10, 6, 9, 0, 12, 11, 7, 13, 15, 1, 3, 14, 5, 2, 8, 4",
			"3, 15, 0, 6, 10, 1, 13, 8, 9, 4, 5, 11, 12, 7, 2, 14", "", "", "", ""},
		// S-Box 5
		{
			"2, 12, 4, 1, 7, 10, 11, 6, 8, 5, 3, 15, 13, 0, 14, 9",
			"14, 11, 2, 12, 4, 7, 13, 1, 5, 0, 15, 10, 3, 9, 8, 6",
			"4, 2, 1, 11, 10, 13, 7, 8, 15, 9, 12, 5, 6, 3, 0, 14",
			"11, 8, 12, 7, 1, 14, 2, 13, 6, 15, 0, 9, 10, 4, 5, 3", "", "", "", ""},
		// S-Box 6
		{
			"12, 1, 10, 15, 9, 2, 6, 8, 0, 13, 3, 4, 14, 7, 5, 11",
			"10, 15, 4, 2, 7, 12, 9, 5, 6, 1, 13, 14, 0, 11, 3, 8",
			"9, 14, 15, 5, 2, 8, 12, 3, 7, 0, 4, 10, 1, 13, 11, 6",
			"4, 3, 2, 12, 9, 5, 15, 10, 11, 14, 1, 7, 6, 0, 8, 13", "", "", "", ""},
		// S-Box 7
		{
			"4, 11, 2, 14, 15, 0, 8, 13, 3, 12, 9, 7, 5, 10, 6, 1",
			"13, 0, 11, 7, 4, 9, 1, 10, 14, 3, 5, 12, 2, 15, 8, 6",
			"1, 4, 11, 13, 12, 3, 7, 14, 10, 15, 6, 8, 0, 5, 9, 2",
			"6, 11, 13, 8, 1, 4, 10, 7, 9, 5, 0, 15, 14, 2, 3, 12", "", "", "", ""},
		// S-Box 8
		{
			"13, 2, 8, 4, 6, 15, 11, 1, 10, 9, 3, 14, 5, 0, 12, 7",
			"1, 15, 13, 8, 10, 3, 7, 4, 12, 5, 6, 11, 0, 14, 9, 2",
			"7, 11, 4, 1, 9, 12, 14, 2, 0, 6, 10, 13, 15, 3, 5, 8",
			"2, 1, 14, 7, 4, 10, 8, 13, 15, 12, 9, 0, 3, 5, 6, 11", "", "", "", ""}}

	stepLabel := widget.NewLabel("Шаг: 0")
	descriptionLabel := widget.NewLabel("Описание текущего шага")
	spaceLabel := widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{Italic: true})
	was2Label := widget.NewLabel("")
	wasLabel := widget.NewLabel("Прошлое состояние данных: ")
	nowLabel := widget.NewLabel("Текущее состояние данных данных: ")
	tableLabel := widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{})
	tableLabels := make([]*widget.Label, 0)
	for i := 0; i < 9; i++ {
		tableLabels = append(tableLabels, widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{}))
	}
	tableLabels[0].Alignment = fyne.TextAlignLeading
	stateLabel := widget.NewLabel("Состояние: ")
	stateLabels := []*widget.Label{stateLabel}
	for i := 0; i < 5; i++ {
		stateLabels = append(stateLabels, widget.NewLabel(""))
	}

	sBoxNextButton := widget.NewButton("Следующая подстановка", func() {
		if currentStep[2] < 8 {
			currentStep[2]++
			updateStepDisplay()
		}

	})

	sBoxPrevButton := widget.NewButton("Предыдущая подстановка", func() {
		if currentStep[2] > 0 {
			currentStep[2]--
			updateStepDisplay()
		}
	})

	fNextButton := widget.NewButton("Следующее действие Фейстеля", func() {
		if currentStep[2] != -1 {
			currentStep[2] = -1
			sBoxNextButton.Disable()
			sBoxPrevButton.Disable()
			updateStepDisplay()
		} else if currentStep[1] < 4 {
			currentStep[1]++
			updateStepDisplay()
		}
		if currentStep[1] == 3 {
			sBoxNextButton.Enable()
			sBoxPrevButton.Enable()
		} else {
			sBoxNextButton.Disable()
			sBoxPrevButton.Disable()
		}
	})

	fPrevButton := widget.NewButton("Предыдущее действие Фейстеля", func() {
		if currentStep[2] != -1 {
			currentStep[2] = -1
			updateStepDisplay()
		} else if currentStep[1] > 0 {
			currentStep[1]--
			updateStepDisplay()
		}
		if currentStep[1] == 3 {
			sBoxNextButton.Enable()
			sBoxPrevButton.Enable()
		} else {
			sBoxNextButton.Disable()
			sBoxPrevButton.Disable()
		}
	})
	rNextButton := widget.NewButton("Следующий раунд", func() {
		fNextButton.Enable()
		fPrevButton.Enable()
		if currentStep[1] != -1 {
			currentStep[2] = -1
			currentStep[1] = -1
			updateStepDisplay()
		} else if currentStep[3] < 15 {
			currentStep[3]++
			updateStepDisplay()
		}
	})

	rPrevButton := widget.NewButton("Предыдущий раунд", func() {
		fNextButton.Enable()
		fPrevButton.Enable()
		if currentStep[1] != -1 {
			currentStep[2] = -1
			currentStep[1] = -1
			updateStepDisplay()
		} else if currentStep[3] > 0 {
			currentStep[3]--
			updateStepDisplay()
		}
	})

	nextButton := widget.NewButton("Следующий шаг", func() {
		if currentStep[3] != -1 {
			currentStep[3] = -1
			currentStep[2] = -1
			currentStep[1] = -1
			updateStepDisplay()
		} else if currentStep[0] < 4 {
			currentStep[0]++
			updateStepDisplay()
		}
		if currentStep[0] == 3 {
			rNextButton.Enable()
			rPrevButton.Enable()
		} else {
			sBoxNextButton.Disable()
			fNextButton.Disable()
			sBoxPrevButton.Disable()
			fPrevButton.Disable()
			rNextButton.Disable()
			rPrevButton.Disable()
		}
		currentStep[3] = -1
		currentStep[2] = -1
		currentStep[1] = -1
	})

	prevButton := widget.NewButton("Предыдущий шаг", func() {
		if currentStep[3] != -1 {
			currentStep[3] = -1
			currentStep[2] = -1
			currentStep[1] = -1
			updateStepDisplay()
		} else if currentStep[0] > 0 {
			currentStep[0]--
			updateStepDisplay()
		}
		if currentStep[0] == 3 {
			rNextButton.Enable()
			rPrevButton.Enable()
		} else {
			sBoxNextButton.Disable()
			fNextButton.Disable()
			sBoxPrevButton.Disable()
			fPrevButton.Disable()
			rNextButton.Disable()
			rPrevButton.Disable()
		}
		currentStep[3] = -1
		currentStep[2] = -1
		currentStep[1] = -1
	})

	returnButton := widget.NewButton("Вернуться на главный экран", func() {
		MainMenu(w)
	})

	// Кнопка для запуска алгоритма

	// Контейнер

	subStepButtons := container.NewHBox(rPrevButton, rNextButton, fPrevButton, fNextButton)

	content = container.NewVBox(
		stepLabel,
		descriptionLabel,
		was2Label,
		wasLabel,
		nowLabel,
		tableLabel,
		spaceLabel,
	)
	for i := 0; i < 8; i++ {
		content.Add(tableLabels[i])
	}
	content.Add(container.NewHBox(sBoxPrevButton, sBoxNextButton))
	content.Add(subStepButtons)
	content.Add(container.NewHBox(prevButton, nextButton, returnButton))

	sBoxNextButton.Disable()
	fNextButton.Disable()
	sBoxPrevButton.Disable()
	fPrevButton.Disable()
	rNextButton.Disable()
	rPrevButton.Disable()

	currentStep[0] = -1
	currentStep[1] = -1
	currentStep[2] = -1
	currentStep[3] = -1

	generateSteps(mainText, mainKey)

	w.SetContent(content)
	w.Resize(fyne.NewSize(800, 600))
	w.CenterOnScreen()
}

func updateStepDisplay() {
	i := currentStep[3]
	j := currentStep[1]
	k := currentStep[2]
	n := currentStep[0]
	o := content.Objects
	o[7].(*widget.Label).Alignment = fyne.TextAlignCenter
	o[0].(*widget.Label).SetText("Шаг: " + strconv.Itoa(n+1) + "." + strconv.Itoa(i+1) + "." + strconv.Itoa(j+1) + "." + strconv.Itoa(k+1))
	if k != -1 {
		if k != 8 {
			for i := 0; i < 8; i++ {
				o[6+i].(*widget.Label).SetText("")
			}
			o[1].(*widget.Label).SetText(helpers[1][0])
			o[2].(*widget.Label).SetText("row = block[0]*2+block[5]. col = block[1]*8+block[2]*4+block[3]*2+block[4]")
			temp := make([]int, 6)
			for a := 0; a < 6; a++ {
				temp[a] = rawStepData[i][3][k][a]
			}
			o[3].(*widget.Label).SetText("row=" + fmt.Sprint(temp[0]) + "*2 + " + fmt.Sprint(temp[5]) + " =" + fmt.Sprint(rawStepData[i][3][k][10]) + "; col=" + fmt.Sprint(temp[1]) + "*8 + " + fmt.Sprint(temp[2]) + "*4 + " + fmt.Sprint(temp[3]) + "*2 + " + fmt.Sprint(temp[4]) + " =" + fmt.Sprint(rawStepData[i][3][k][11]))
			o[4].(*widget.Label).SetText("Блок:" + fmt.Sprint(rawStepData[i][2][0][k:k*6]) + "->" + fmt.Sprint(rawStepData[i][2][0][k*6:k*6+6]) + "<-" + fmt.Sprint(rawStepData[i][2][0][k*6+6:]))
			o[5].(*widget.Label).SetText("Подблок №" + fmt.Sprint(k) + ": " + fmt.Sprint(rawStepData[i][3][k][:6]))
			o[6].(*widget.Label).SetText("Полученный подблок" + fmt.Sprint(rawStepData[i][3][k][6:10]) + "=" + fmt.Sprint(rawStepData[i][3][k][12]))
			indexL := strings.Index(tables[8+k][rawStepData[i][3][k][10]], fmt.Sprint(rawStepData[i][3][k][12]))
			o[7].(*widget.Label).Alignment = fyne.TextAlignLeading
			o[7].(*widget.Label).SetText("S-box №" + fmt.Sprint(k) + ":")
			for a := 0; a < 4; a++ {
				st := tables[8+k][a][:indexL] + "->"
				st += tables[8+k][a][indexL : indexL+rawStepData[i][3][k][12]/10+1]
				st += "<-" + tables[8+k][a][indexL+rawStepData[i][3][k][12]/10+1:]
				if a == rawStepData[i][3][k][10] {
					o[8+a].(*widget.Label).SetText(st)
				} else {
					o[8+a].(*widget.Label).SetText(tables[8+k][a])
				}

			}
		} else {
			o[1].(*widget.Label).SetText(helpers[1][1])
			o[2].(*widget.Label).SetText("") //PT
			o[3].(*widget.Label).SetText("")
			o[4].(*widget.Label).SetText("Блок:" + fmt.Sprint(rawStepData[i][3][0][13:45]))
			o[5].(*widget.Label).SetText("Блок:" + fmt.Sprint(rawStepData[i][3][0][45:77]))
			o[6].(*widget.Label).SetText("Таблица перестановки:")
			for i := 0; i < 8; i++ {
				o[7+i].(*widget.Label).SetText(tables[3][i])
			}
		}
	} else if j != -1 {
		o[1].(*widget.Label).SetText(helpers[0][j])
		if j == 0 {
			o[2].(*widget.Label).SetText("") // left | right
			o[3].(*widget.Label).SetText("")
			o[4].(*widget.Label).SetText("Блок:" + fmt.Sprint(rawStepData[0][j][0]))
			o[5].(*widget.Label).SetText("Блок:" + fmt.Sprint(rawStepData[0][j][0][:32]) + " | " + fmt.Sprint(rawStepData[0][j][0][32:]))
			for i := 0; i < 9; i++ {
				o[6+i].(*widget.Label).SetText("")
			}
		} else if j == 1 {
			o[2].(*widget.Label).SetText("") //expand right
			o[3].(*widget.Label).SetText("")
			o[4].(*widget.Label).SetText("Блок:" + fmt.Sprint(rawStepData[0][j-1][0][32:]))
			o[5].(*widget.Label).SetText("Блок:" + fmt.Sprint(rawStepData[0][j][0]))
			o[6].(*widget.Label).SetText("Таблица расширения:")
			for i := 0; i < 8; i++ {
				o[7+i].(*widget.Label).SetText(tables[2][i])
			}
		} else if j == 2 {
			o[2].(*widget.Label).SetText("") //xor
			o[4].(*widget.Label).SetText("Ключ: " + fmt.Sprint(rawKeysData[0][1]))
			o[3].(*widget.Label).SetText("Блок: " + fmt.Sprint(rawStepData[0][j-1][0]))
			o[5].(*widget.Label).SetText("Блок: " + fmt.Sprint(e.Xor(rawKeysData[0][1], rawStepData[0][j-1][0])))
			for i := 0; i < 9; i++ {
				o[6+i].(*widget.Label).SetText("")
			}
		} else if j == 3 {
			o[2].(*widget.Label).SetText("")
			o[3].(*widget.Label).SetText("")
			o[4].(*widget.Label).SetText("Блок:" + fmt.Sprint(rawStepData[0][j-1][0]))
			o[5].(*widget.Label).SetText("Блок:" + fmt.Sprint(rawStepData[0][j][0][45:]))
			for i := 0; i < 9; i++ {
				o[6+i].(*widget.Label).SetText("")
			}
		} else if j == 4 {
			o[2].(*widget.Label).SetText("")
			o[3].(*widget.Label).SetText("")
			o[4].(*widget.Label).SetText("Блок:" + fmt.Sprint(rawStepData[0][0][0][32:]) + " | " + fmt.Sprint(rawStepData[0][j][0][:32]))
			o[5].(*widget.Label).SetText("Блок:" + fmt.Sprint(rawStepData[0][j][0][:32]) + " | " + fmt.Sprint(rawStepData[0][0][0][32:]))
			for i := 0; i < 9; i++ {
				o[6+i].(*widget.Label).SetText("")
			}

		}
	} else if i != -1 {
		if i == 0 {
			o[1].(*widget.Label).SetText(helpers[3][2] + ". Раунд: " + strconv.Itoa(i))
			o[2].(*widget.Label).SetText("") //
			o[3].(*widget.Label).SetText("Выполняется раунд сети Фейстеля")
			o[4].(*widget.Label).SetText("Блок:" + fmt.Sprint(rawStepData[16][3][0]))
			o[5].(*widget.Label).SetText("Блок:" + fmt.Sprint(rawStepData[i][5][0]))
			for i := 0; i < 9; i++ {
				o[6+i].(*widget.Label).SetText("")
			}
		} else {
			o[1].(*widget.Label).SetText(helpers[3][2] + ". Раунд: " + strconv.Itoa(i))
			o[2].(*widget.Label).SetText("") //
			o[3].(*widget.Label).SetText("Выполняется раунд сети Фейстеля")
			o[4].(*widget.Label).SetText("Блок:" + fmt.Sprint(rawStepData[i-1][5][0]))
			o[5].(*widget.Label).SetText("Блок:" + fmt.Sprint(rawStepData[i][5][0]))
			for i := 0; i < 9; i++ {
				o[6+i].(*widget.Label).SetText("")
			}
		}
	} else {
		if n == 0 {
			o[1].(*widget.Label).SetText(helpers[3][1])
			o[2].(*widget.Label).SetText("") //Padding
			o[3].(*widget.Label).SetText("Блок заполняется байтами, равными числу недостающих байт")
			o[4].(*widget.Label).SetText("Блок:" + fmt.Sprint(rawStepData[16][0][0]))
			o[5].(*widget.Label).SetText("Блок:" + fmt.Sprint(rawStepData[16][1][0]))
			for i := 0; i < 9; i++ {
				o[6+i].(*widget.Label).SetText("")
			}
		}
		if n == 1 {
			o[1].(*widget.Label).SetText(helpers[3][0])
			o[2].(*widget.Label).SetText("") //IP
			o[3].(*widget.Label).SetText("Начальная перестановка")
			o[4].(*widget.Label).SetText("Исходный блок:" + fmt.Sprint(rawStepData[16][1][0]))
			o[5].(*widget.Label).SetText("Выходной блок:" + fmt.Sprint(rawStepData[16][2][0]))
			o[6].(*widget.Label).SetText("Таблица начальной перестановки:")
			for i := 0; i < 8; i++ {
				o[7+i].(*widget.Label).SetText(tables[1][i])
			}
		}
		if n == 2 {
			o[1].(*widget.Label).SetText(helpers[3][4])
			o[2].(*widget.Label).SetText("") //
			o[3].(*widget.Label).SetText("")
			o[4].(*widget.Label).SetText("Основной Ключ:" + fmt.Sprint(rawStepData[16][3][0]))
			o[5].(*widget.Label).SetText("")
			for i := 0; i < 9; i++ {
				o[6+i].(*widget.Label).SetText("")
			}
		}
		if n == 3 {
			o[1].(*widget.Label).SetText(helpers[3][2])
			o[2].(*widget.Label).SetText("") // feistel
			o[3].(*widget.Label).SetText("")
			o[4].(*widget.Label).SetText("Исходный Блок:" + fmt.Sprint(rawStepData[16][2][0]))
			o[5].(*widget.Label).SetText("Выходной Блок:" + fmt.Sprint(rawStepData[15][5][0]))
			for i := 0; i < 9; i++ {
				o[6+i].(*widget.Label).SetText("")
			}
		}
		if n == 4 {
			o[1].(*widget.Label).SetText(helpers[3][3])
			o[2].(*widget.Label).SetText("") //FP
			o[3].(*widget.Label).SetText("")
			o[4].(*widget.Label).SetText("Исходный Блок:" + fmt.Sprint(rawStepData[15][5][0]))
			o[5].(*widget.Label).SetText("Выходной Блок:" + fmt.Sprint(rawStepData[16][4][0]))
			o[6].(*widget.Label).SetText("Таблица конечной перестановки:")
			for i := 0; i < 8; i++ {
				o[7+i].(*widget.Label).SetText(tables[6][i])
			}
		}
	}
}

var rawStepData [17][6][][]int // 17 - initials
var rawKeysData [17][2][]int   //17 - PC1

func generateSteps(t, k string) {
	text := e.ByteSToBinS([]byte(t))
	if len(text) > 64 {
		text = text[:64]
	}
	rawStepData[16][0] = append(rawStepData[16][0], text)
	text = e.ByteSToBinS(e.AddDESPadding(e.BinSToByteS(text)))
	rawStepData[16][1] = append(rawStepData[16][1], text)
	text = e.Permute(text, e.IPTable)
	rawStepData[16][2] = append(rawStepData[16][2], text)
	key := e.ByteSToBinS([]byte(k))[:64]
	rawStepData[16][3] = append(rawStepData[16][3], key)

	keys := GenerateKeys(key)

	for i := 0; i < 16; i++ {
		left := text[:32]
		right := text[32:]

		rawStepData[i][0] = append(rawStepData[i][0], text)

		right = e.Expand(right)
		rawStepData[i][1] = append(rawStepData[i][1], right)

		right = e.Xor(right, keys[i])
		rawStepData[i][2] = append(rawStepData[i][2], right)

		right = generateSboxSteps(right, i)

		rawStepData[i][4] = append(rawStepData[i][4], right)

		text = append(right, left...)
		rawStepData[i][5] = append(rawStepData[i][5], text)
	}
	text = e.Permute(text, e.FPTable)
	rawStepData[16][4] = append(rawStepData[16][4], text)
}

func generateSboxSteps(right []int, i int) []int {
	rawStepData[i][3] = make([][]int, 8)
	output := make([]int, 0, 32)
	for j := 0; j < 8; j++ {
		row := right[j*6]*2 + right[j*6+5]
		col := right[j*6+1]*8 + right[j*6+2]*4 + right[j*6+3]*2 + right[j*6+4]
		binary := fmt.Sprintf("%04b", e.SBox[j][row][col])
		for _, bit := range binary {
			output = append(output, int(bit-'0'))
		}
		rawStepData[i][3][j] = append(rawStepData[i][3][j], right[j*6:j*6+6]...)
		rawStepData[i][3][j] = append(rawStepData[i][3][j], output[j*4:j*4+4]...)
		rawStepData[i][3][j] = append(rawStepData[i][3][j], row, col, e.SBox[j][row][col])
	}
	rawStepData[i][3][0] = append(rawStepData[i][3][0], output...)
	output = e.Permute(output, e.PermutationTable)
	rawStepData[i][3][0] = append(rawStepData[i][3][0], output...)
	return output
}

func GenerateKeys(key []int) [][]int {
	key = e.Permute(key, e.PC1)
	rawKeysData[16][0] = key
	keys := make([][]int, 16)
	left := key[:28]
	right := key[28:]
	for i := 0; i < 16; i++ {
		left = e.LeftShift(left, e.LeftShifts[i])
		right = e.LeftShift(right, e.LeftShifts[i])
		rawKeysData[i][0] = append(left, right...)
		key = append(left, right...)
		keys[i] = e.Permute(key, e.PC2)
		rawKeysData[i][1] = append(rawKeysData[i][1], keys[i]...)
	}
	return keys
}
