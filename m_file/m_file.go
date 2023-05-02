// Пакет tools, мой пакет с моими инструментами для упрощения работы.
package m_file

import (
	"bufio"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
)

// Функция GetPathFiles ищет файлы или папки с нужным именем и влзвращает массив путей к ним.
// Папки ищутся точно по именни. А файлы можно найти, частично введя название.
func GetPathFiles(pathFolderSearch string, nameSearch string, IsFile bool) []string {
	result := []string{}
	pathFolderSearch = strings.Replace(pathFolderSearch, "\\", "/", -1)

	resultSearch := searchAll(pathFolderSearch)
	resultPath := searchThisFileDir(resultSearch, nameSearch, IsFile)

	if len(resultSearch) != 0 {
		return resultPath
	}
	//nameSearchFiles = strings.ToLower(nameSearchFiles)

	return result
}

// Функция searchThisFileDir фильрует папки или файлы по нужному слову и возвращает массив путей к ним.
func searchThisFileDir(paths map[string][]string, nameSearch string, IsFile bool) []string {
	nameSearch = strings.ToLower(nameSearch)
	result := []string{}

	if IsFile {
		for _, file := range paths["files"] {
			nameFile := strings.ToLower(path.Base(file))
			if strings.Contains(nameFile, nameSearch) {
				result = append(result, file)
			}
		}

	} else {
		for _, file := range paths["dirs"] {
			nameFile := strings.ToLower(path.Base(file))
			if nameSearch == nameFile {
				result = append(result, file)
			}
		}

	}

	return result
}

// Фнкция searchDir ищет все папки и файлы в пути и возвращает в виде map[string][]string{"dirs": {}, "files": {}}
func searchAll(pathDir string) map[string][]string {
	resultSearch := map[string][]string{"dirs": {}, "files": {}}

	folders, err := os.ReadDir(pathDir)
	if err != nil {
		log.Fatal(fmt.Printf("Err: При сканировании папки \n%v\n", err))
	}

	dirArr := []string{}
	for _, name := range folders {
		p := path.Join(pathDir, name.Name())

		if name.IsDir() {
			dirArr = append(dirArr, p)
			resultSearch["dirs"] = append(resultSearch["dirs"], p)

		} else {
			resultSearch["files"] = append(resultSearch["files"], p)
		}
	}

	if len(dirArr) != 0 {
		for _, d := range dirArr {
			resultNext := searchAll(d)
			for key, value := range resultNext {
				resultSearch[key] = append(resultSearch[key], value...)
			}
		}
	}

	return resultSearch
}

// Функция нужна для ввода данных с клавиатуры
func InputUser(message string) string {
	fmt.Printf("%v", message)
	scan := bufio.NewReader(os.Stdin)
	i, err := scan.ReadString('\n')
	if err != nil {
		panic(err)
	}
	i = strings.TrimSpace(i)
	return i
}

// Функция считивает файл и возвращает массив строк
func File_open_In_arr(path_cookie string) ([]string, error) {
	file, err := os.OpenFile(path_cookie, os.O_RDONLY, 0777)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	sc := bufio.NewScanner(file)
	var arr_str []string
	for sc.Scan() {
		str := strings.TrimSpace(sc.Text())
		if str != "" {
			arr_str = append(arr_str, str)
		}

	}
	return arr_str, nil
}

// Функция читает папку и возвращает []fs.DirEntry
func Read_dir(path_dir string) ([]fs.DirEntry, error) {
	zip_folder, err := os.ReadDir(path_dir)
	if err != nil {
		return []fs.DirEntry{}, err
	}
	return zip_folder, nil
}

// Функция читает файл и возвращает ответ в байтах
func File_open(path_file string) ([]byte, error) {
	file, err := ioutil.ReadFile(path_file)
	if err != nil {
		return []byte{}, err
	}
	return file, nil
}
