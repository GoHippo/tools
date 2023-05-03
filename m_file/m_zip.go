package m_file

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
)

// Функция extract_zip( path_zip - путь к zip архиву, path_folder_extract - путь к папку для распаковки,
// если nil, то распакует по имени архива ) распаковывает zip архива.
func Zip_Extract(path_zip string, path_dir_extract string) error {
	//path_zip := `temp_pac/example.zip`

	if path_dir_extract == "" {
		path_dir_extract = path.Join(path.Dir(path_zip), path.Base(path_zip)[:len(path.Base(path_zip))-4])
	}

	if path_dir_extract[:1] == "/" {
		path_dir_extract = "." + path_dir_extract
	}

	// Открываем zip архив
	z, err := zip.OpenReader(path_zip)

	if err != nil {
		return fmt.Errorf("Ошибка при открытии архива:%v", err)
	}
	defer z.Close()

	// Создаем директорию, удаляя старую, если она существует.
	os.RemoveAll(path_dir_extract)
	err = os.MkdirAll(path_dir_extract, 0755)
	if err != nil {
		return fmt.Errorf("Ошибка при создании директории:%v Err:%v", path_dir_extract, err)
	}

	// Распаковываем файлы
	for _, f := range z.File {

		// Получаем полный путь к файлу
		path_file := filepath.Join(path_dir_extract, f.Name)

		// Получаем путь к директориям файла
		var path_dir = filepath.ToSlash(path_file)

		// Создаем директории для файла, иначе записываем файл
		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(path_dir, f.Mode()); err != nil {
				return fmt.Errorf("Ошибка при создании папки:%v Err:%v", path_dir, err)
			}
		} else {

			// Создаем новый файл
			dst, err := os.Create(path_file)
			if err != nil {
				fmt.Printf("Ошибка при создании нового файла:%v\n", err)
				continue
			}
			defer dst.Close()

			// Открываем файл в архиве
			src, err := f.Open()
			if err != nil {
				fmt.Printf("Ошибка при открытии файла в архиве:%v\n", err)
				continue
			}
			defer src.Close()

			// Копируем содержимое файла в новый файл
			_, err = io.Copy(dst, src)
			if err != nil {
				fmt.Printf("Ошибка при распаковке файла:%v\n", err)
				continue
			}
		}
	}
	return nil
}
