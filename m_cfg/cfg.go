package m_cfg

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Setting struct {
	json_cfg        // Чтобы изменить настройки конфиги, нужно только заменить их тут в структуре
	path_cfg string //путь к конфигу
}

// Важно!! Не забыть изменить поле json:"" иначе имя будет сохранятся старое
type json_cfg struct {
	Path_folder_logs      string   `json:"path_folder_logs"`      // путь к папке с логами
	Path_folder_result    string   `json:"path_folder_result"`    // путь к результатам поиска по почтам
	Save_full_logs        bool     `json:"save_full_logs"`        // сохранение полного лога
	Google_mail           []string `json:"google_mail"`           // запросы по ключевым словам gmail
	Outlook_mail          []string `json:"outlook_mail"`          // запросы по ключевым словам outlook mail
	Yahoo_mail            []string `json:"yahoo_mail"`            // запросы по ключевым словам yahoo mail
	Epicgames             bool     `json:"epicgames"`             // чекает эпик куки на вход
	Youtube               bool     `json:"youtube"`               // чекер youtube
	Youtube_min_subscribe int      `json:"youtube_min_subscribe"` // минимальное число подписчиков для сохранения канала в отдельную папку с инфо
	Proxy                 string   `json:"proxy"`                 // прокси в виде строки "127.0.0.1:8888"
	Max_tcp_connections   int      `json:"max_tcp_connections"`   // максимальное количество tcp подключенний
	MaxProc               int      `json:"maxProc"`               // максимум активных процессоров
	Max_threads           int      `json:"max_threads"`           // количество одновременно чекающих логов(папки cookies)
	Logger_on             bool     `json:"logger_on"`             //включение логгера приложения
}

// Метод сохраняет изменения в конфиге
func (c *Setting) SaveCfg() {
	js, err := json.Marshal(c)
	check_err(err, "CreateJson Marshal")

	js_wr := make([]byte, 0)
	n := 0
	for _, k := range js {
		if k == byte('{') {
			js_wr = append(js_wr, '{', '\n', '\t')
			n = +1
			continue
		}
		if k == byte('}') {
			js_wr = append(js_wr, '\n')
		}

		if k == ',' {
			js_wr = append(js_wr, ',', '\n')

			i := 0
			for i != n {
				i += 1
				js_wr = append(js_wr, '\t')
			}
			continue
		}

		js_wr = append(js_wr, k)
	}

	err = os.WriteFile(c.path_cfg, js_wr, 0777)
	check_err(err, "CreateJson OpenFile")
}

// Метод открывает файл и возвращает текс в виде байт
func (cfg *Setting) openJson() ([]byte, error) {
	file, err := os.OpenFile(cfg.path_cfg, os.O_RDONLY, 0777)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	sc := bufio.NewScanner(file)
	var data []byte
	for sc.Scan() {
		t := sc.Text()

		data = append(data, t...)
	}
	return data, nil
}

// Метод проверяет конфиг(json) на разные ошибки
func (cfg *Setting) checkJsonCfg() {

	js, err := cfg.openJson()
	if err != nil {
		if strings.Contains(err.Error(), "The system cannot find the file specified") {
			fmt.Printf("Файл по пути найстроек не найден, создаю новый!\n")
			cfg.SaveCfg()
			return
		} else {
			check_err(err, "checkJsonCfg OpenFile")
		}
	}

	var data *map[string]interface{}

	err = json.Unmarshal(js, &data)
	if err != nil {
		fmt.Printf("Ошибка в файле пересозданно на дефолт. err: %v\n", err)
		cfg.SaveCfg()
		return
	}

	var check_key *map[string]interface{}
	var s Setting
	j, _ := json.Marshal(s)
	err = json.Unmarshal(j, &check_key)
	check_err(err, "checkJsonCfg Unmarshal")
	//fmt.Println(check_key)
	for key, v := range *check_key {
		_ = v
		_, ok := (*data)[key]
		if !ok {
			fmt.Println("Файл настроек поврежден. пересозданно на дефолт.")
			cfg.SaveCfg()
			return
		}
	}
}

// Функция получает конфиг из указанного пути
func Get_cfg(path_json string) *Setting {
	var cfg *Setting = &Setting{path_cfg: path_json}

	cfg.checkJsonCfg()

	js, err := cfg.openJson()
	check_err(err, "main openJson")

	err = json.Unmarshal(js, &cfg)
	check_err(err, "main unmarshal")

	return cfg
}

// Функция проверяет ошибки
func check_err(err error, f string) {
	if err != nil {
		panic(fmt.Errorf("[%v] %v", f, err))
	}
}
