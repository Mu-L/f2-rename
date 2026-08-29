package rename

import (
	"errors"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/ayoisaiah/f2/v2/internal/file"
	"github.com/ayoisaiah/f2/v2/internal/osutil"
)

func createBackupFile(fileName string) (*os.File, error) {
	backupFilePath := filepath.Join(
		os.TempDir(),
		"f2",
		"backups",
		fileName,
	)

	slog.Debug(
		"creating backup file",
		slog.String("backup_file", backupFilePath),
	)

	err := os.MkdirAll(filepath.Dir(backupFilePath), osutil.DirPermission)
	if err != nil {
		return nil, err
	}

	// Create or truncate backupFile
	return os.Create(backupFilePath)
}

// backupChanges records the details of a renaming operation to the specified
// writer so that it may be reverted if necessary. If a writer is not specified
// it records the changes to the filesystem.
func backupChanges(
	changes file.Changes,
	cleanedDirs []string,
	fileName string,
	w io.Writer,
) (err error) {
	if w == nil {
		backupFile, err := createBackupFile(fileName)
		if err != nil {
			return err
		}

		w = backupFile

		defer func() {
			err = errors.Join(err, backupFile.Close())
		}()
	}

	b := file.Backup{
		Changes:     changes,
		CleanedDirs: cleanedDirs,
	}

	slog.Debug("backing up changed", slog.Any("backup", b))

	return b.RenderJSON(w)
}
