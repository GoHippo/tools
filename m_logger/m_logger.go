// Пакет содержит утилиту logger. Чтобы воспользоваться нужно обьявить обьектc параметрами. После закрыть.
//
//		l:=NewLogs(path_save_file string--путь сохранения, timer_save_logs time.Duration--время записи, print_log bool--вывод на экран)
//	 defer l.Close() -- обязательно закрыть логгер
package m_logger

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Logs struct {
	str           string    //Архив добавленных логов-строк.
	path_save     string    //путь к сохранению лога
	save_time_log int       //время для сохранения и очищения архива строк в фойл
	channel_exit  chan bool // Канал для закрытия горутины
	//channel_connect chan string	//	Канал для предачи лога из другого логгера
	print_log bool //печатает лог в терминал
	On        bool // выключает логгер если нужно
	save_file bool // включить сохранения лога в файл
}

// Флаги для обозначения в файле
const L_warn = "[Warn]"
const L_err = "[Error]"
const L_info = "[Info]"
const L_fatal = "[Fatal]"

// мутекс для блокировки общего доступа во время записи в архив логов и удаления
var mutex sync.Mutex

func (l *Logs) DeleteFileLogs() {
	l.save_file = false
	_, err := os.ReadFile(l.path_save)
	if err == nil {
		os.Remove(l.path_save)
	}
}

// Метод управляет печатьб лога на экран
func (l *Logs) SetPrint(print_logs bool) {
	l.print_log = print_logs
}

// Метод сохраняет лог по пути в структуре
// Очишает скопившийся архив лог-строк в структуре "Str"
func (l *Logs) save_log() {
	if l.str == "" {
		return
	}
	mutex.Lock()
	log_new := l.str
	l.str = ""
	mutex.Unlock()

	if l.save_file {
		file, err := os.OpenFile(l.path_save, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0777)
		if err != nil {
			log.Fatal("[Err]Файл лога не открывается для записи! err: ", err)
		}
		defer file.Close()

		if _, err = file.Write([]byte(log_new)); err != nil {
			log.Fatal("[Err] Ошибка при записи в открый файл лога! err: ", err)
		}
	}
}

// Метод добавляет новый лог в архив
// Указать строку и флаг. например: log.Append("File create!",tools.L_log)
func (l *Logs) L_Append(flag string, s ...any) {
	if !l.On {
		return
	}

	l.check_log_nil()

	// var string_s string
	// for _,i:=range s{
	// 	string_s+=fmt.Sprintf("dsd",s...)
	// }
	stroka := fmt.Sprintf("%v", s...)

	t := time.Now().Format("2006-01-02T15:04:05")

	gorutine_id := goid()

	new_log := fmt.Sprintf("\n%v[%v][%v] \t%v", flag, t, gorutine_id, stroka[1:len(stroka)-1])

	mutex.Lock()
	l.str = l.str + new_log
	mutex.Unlock()
	if l.print_log {
		fmt.Println(new_log)
	}
}

// Функция добавляет лог с флагом [Info]
func (l *Logs) L_Info(s ...any) {
	l.L_Append(L_info, s)
}

// Функция добавляет лог с флагом [Err]
func (l *Logs) L_Err(s ...any) {
	l.L_Append(L_err, s)
}

// Функция добавляет лог с флагом [Fatal]
func (l *Logs) L_Fatal(s ...any) {
	l.L_Append(L_fatal, s)
}

// Функция добавляет лог с флагом [Warn]
func (l *Logs) L_Warning(s ...any) {
	l.L_Append(L_warn, s)
}

// Запускает логирвоание в режиме постоянной проверки архива логов в обьекте и сохранение через время
func (l *Logs) start() {
	l.check_log_nil()

	i := 0
	for {
		select {
		case <-l.channel_exit:
			return
		// case <-l.channel_connect:
		// 	l.str=l.str+<-l.channel_connect
		default:
			time.Sleep(time.Second * 1) //l.save_time_log)
		}

		i += 1
		if i%l.save_time_log == 0 {
			l.save_log()
		}
	}
}

// Метод StopWrite() останавливает логирование и сохраняет последние добавленные логи.
func (l *Logs) Close() {
	if !l.On {
		return
	}
	l.check_log_nil()

	err := recover()
	pan, err := rec(l, err)

	l.L_Info("Logging stop...")
	l.save_log()

	l.channel_exit <- true
	if pan {
		panic(err)
	}
}

// Функция NewLogs создает новый обьект для логгирования,нужно указать путь к сохранению лога.
// Создает новый файл лога по пути.
//
//	log:=NewLogs("./logs.txt",30,true)  timer_save_logs in second
func NewLogs(path_save_file string, timer_save_logs int, print_log bool) *Logs {
	if path_save_file == "" {
		path_save_file = "/logs.txt"
	}
	if timer_save_logs == 0 {
		timer_save_logs = 1
	}

	ch_exit := make(chan bool)
	//ch_connect=make(chan string)
	log_object := &Logs{str: "---Start Logs---", path_save: path_save_file, save_time_log: timer_save_logs, print_log: true, channel_exit: ch_exit, On: true, save_file: true} //,channel_connect: ch_connect}

	if log_object.save_file {
		if err := os.WriteFile(path_save_file, []byte(""), 0777); err != nil {
			log.Fatal("[err] Ошибка в создании  лог файла! err:", err)
		}
	}

	go log_object.start()

	return log_object
}

// Функция rec() Проверяет наличие ошибки и записывает ее в лог,если она есть.
// Возвращает bool и ошибку, чтобы вызвать панику дальше,если она была.
func rec(l *Logs, err interface{}) (bool, interface{}) {
	l.check_log_nil()

	if err != nil {
		stack := string(debug.Stack())

		switch x := err.(type) {
		case error:
			l.L_Fatal(x.Error(), L_fatal)
			l.L_Append("[Stack_Fatal]", stack)
		case string:
			l.L_Fatal(x, L_fatal)
			l.L_Append("[Stack_Fatal]", stack)
		}
		return true, err
	}
	return false, nil
}

func (l *Logs) check_log_nil() {
	if l == nil {
		panic("Обьект логгера ранвен nil")
	}
}

func goid() int {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	id, err := strconv.Atoi(idField)
	if err != nil {
		panic(fmt.Sprintf("cannot get goroutine id: %v", err))
	}
	return id
}
