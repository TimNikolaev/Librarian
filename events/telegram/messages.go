package telegram

const msgHelp = `Reference:

I can save and keep you pages. Also I can offer you them to read.

In order to save the page, just send me al link to it.

I order to get a random page from your list,  send me command /rnd.

Caution! After that, this page will be removed from your list!

Справка:

Я могу сохранить ваши страницы. Также я могу предложить вам их для прочтения.

Чтобы сохранить страницу, просто пришлите мне ссылку на нее.

Чтобы выбрать случайную страницу из вашего списка, отправьте мне команду /rnd.

Осторожно! После этого эта страница будет удалена из вашего списка!`

const msgHello = "Hi there! \n\n" + msgHelp

const (
	msgUnknownCommand = "Unknown command 🧐"
	msgNoSavedPages   = "You have no saved pages 🙊"
	msgSaved          = "Saved! 👌"
	msgAlreadyExists  = "You have already have this page in your list 😊"
)
