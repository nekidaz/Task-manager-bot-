package bot

import (
	"gopkg.in/telebot.v3"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"
	"upgrade/internal/models"
)

type UpgradeBot struct {
	Bot   *telebot.Bot
	Users *models.UserModel
	Tasks *models.TaskModel
}

var gameItems = [3]string{
	"камень",
	"ножницы",
	"бумага",
}

var winSticker = &telebot.Sticker{
	File: telebot.File{
		FileID: "CAACAgIAAxkBAAEGMEZjVspD4JulorxoH7nIwco5PGoCsAACJwADr8ZRGpVmnh4Ye-0RKgQ",
	},
	Width:    512,
	Height:   512,
	Animated: true,
}

var loseSticker = &telebot.Sticker{
	File: telebot.File{
		FileID: "CAACAgIAAxkBAAEGMEhjVsqoRriJRO_d-hrqguHNlLyLvQACogADFkJrCuweM-Hw5ackKgQ",
	},
	Width:    512,
	Height:   512,
	Animated: true,
}

func (bot *UpgradeBot) StartHandler(ctx telebot.Context) error {
	newUser := models.User{
		Name:       ctx.Sender().Username,
		TelegramId: ctx.Chat().ID,
		FirstName:  ctx.Sender().FirstName,
		LastName:   ctx.Sender().LastName,
		ChatId:     ctx.Chat().ID,
	}

	existUser, err := bot.Users.FindOne(ctx.Chat().ID)

	if err != nil {
		log.Printf("Ошибка получения пользователя %v", err)
	}

	if existUser == nil {
		err := bot.Users.Create(newUser)

		if err != nil {
			log.Printf("Ошибка создания пользователя %v", err)
		}
	}

	return ctx.Send("Привет, " + ctx.Sender().FirstName)
}

func (bot *UpgradeBot) GameHandler(ctx telebot.Context) error {
	return ctx.Send("Сыграем в камень-ножницы-бумага " +
		"Введи твой вариант в формате /try камень")
}

func (bot *UpgradeBot) TryHandler(ctx telebot.Context) error {
	attempts := ctx.Args()

	if len(attempts) == 0 {
		return ctx.Send("Вы не ввели ваш вариант")
	}

	if len(attempts) > 1 {
		return ctx.Send("Вы ввели больше одного варианта")
	}

	try := strings.ToLower(attempts[0])
	botTry := gameItems[rand.Intn(len(gameItems))]

	if botTry == "камень" {
		switch try {
		case "ножницы":
			ctx.Send(loseSticker)
			ctx.Send("🪨")
			return ctx.Send("Камень! Ты проиграл!")
		case "бумага":
			ctx.Send(winSticker)
			ctx.Send("🪨")
			return ctx.Send("Камень! Ты выиграл!")
		}
	}

	if botTry == "ножницы" {
		switch try {
		case "камень":
			ctx.Send(winSticker)
			ctx.Send("✂️")
			return ctx.Send("Ножницы! Ты выиграл!")
		case "бумага":
			ctx.Send(loseSticker)
			ctx.Send("✂️")
			return ctx.Send("Ножницы! Ты проиграл!")
		}
	}

	if botTry == "бумага" {
		switch try {
		case "ножницы":
			ctx.Send(winSticker)
			ctx.Send("📃")
			return ctx.Send("Бумага! Ты выиграл!")
		case "камень":
			ctx.Send(loseSticker)
			ctx.Send("📃")
			return ctx.Send("Бумага! Ты проиграл!")
		}
	}

	if botTry == try {
		return ctx.Send("Ничья!")
	}

	return ctx.Send("Кажется вы ввели неверный вариант!")
}

func InitBot(token string) *telebot.Bot {
	pref := telebot.Settings{
		Token:  token,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := telebot.NewBot(pref)

	if err != nil {
		log.Fatalf("Ошибка при инициализации бота %v", err)
	}

	return b
}

func (bot *UpgradeBot) AddTask(ctx telebot.Context) error {

	// принимаю текст
	attems := strings.ReplaceAll(ctx.Text(), "/addTask", "")

	if attems == "" {
		return ctx.Send("Надо ввести значение по как /addTask Заголовк | Текст | дата окончания\t Например\t /addTask Сходить в магаз | Пойти в магаз за хлебом :) | 19:30   ")
	}
	//разделяю на части как заголовк-задача-время
	task := strings.Split(attems, "|")

	if len(task) < 3 || len(task) > 3 {
		return ctx.Send("Надо ввести значение по как /addTask Заголовк | Текст | дата окончания\t Например\t /addTask Сходить в магаз | Пойти в магаз за хлебом :) | 19:30   ")
	}

	title := task[0]
	description := task[1]
	endTime := task[2]

	if title == "" || description == "" || endTime == "" {
		return ctx.Send("Надо ввести значение по как /addTask Заголовк | Текст | дата окончания\t Например\t /addTask Сходить в магаз | Пойти в магаз за хлебом :) | 19:30   ")
	}

	newTask := models.Tasks{
		Title:   title,
		Descr:   description,
		EndDate: endTime,
		Userid:  ctx.Chat().ID,
	}

	err := bot.Tasks.CreateTask(newTask)

	if err != nil {
		log.Printf("Ошибка создания задания %v", err)
	}

	return ctx.Send("как скажешь:)")

}

func (bot *UpgradeBot) ShowTasks(context telebot.Context) error {

	tasksUser := &bot.Tasks
	_, tasksList := (*tasksUser).ShowTaskDb(context.Chat().ID)

	if len(tasksList) == 0 {
		return context.Send("У вас нет заданий ")
	}

	for i := 0; i < len(tasksList); i++ {
		context.Send(strconv.Itoa(i+1) + "." + tasksList[i].Title + ":" + tasksList[i].Descr)
	}

	return context.Send(tasksList)
}

func (bot *UpgradeBot) DeleteTask(context telebot.Context) error {

	tasksUser := &bot.Tasks

	apptemp := context.Args()

	if len(apptemp) > 1 || len(apptemp) < 1 {
		return context.Send("попробуйте еще раз")
	}

	askId, _ := strconv.ParseInt(context.Args()[0], 0, 64)

	if askId < 1 {
		return context.Send("попробуйте еще раз")
	}

	_, taskslist := (*tasksUser).ShowTaskDb(context.Chat().ID)

	if int(askId) > len(taskslist) {
		return context.Send("неправильный ввод")
	}

	if len(taskslist) == 0 {
		return context.Send("У вас нет заданий ")
	}

	(*tasksUser).DeleteTask(taskslist[askId-1].Id)

	return context.Send("Задание : " + taskslist[askId-1].Title + ":  удалено")
}

func (bot *UpgradeBot) CommandsList(context telebot.Context) error {
	return context.Send("Команды:\n" +
		"1./addTasks\n" +
		"2./showTasks\n" +
		"3./deleteTask\n" +
		"3./deleteTask\n" +
		"4./game\n" +
		"5./try")
}
