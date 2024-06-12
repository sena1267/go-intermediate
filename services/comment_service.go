package services

import (
	"github.com/sena1267/go-intermediate/apperrors"
	"github.com/sena1267/go-intermediate/models"
	"github.com/sena1267/go-intermediate/repositories"
)

// PostCommentHandlerで使用することを想定したサービス
// 引数の情報をもとに新しいコメントを作り、結果を返却
func (s *MyAppService) PostCommentService(comment models.Comment) (models.Comment, error) {
	newComment, err := repositories.InsertComment(s.db, comment)
	if err != nil {
		err = apperrors.InsertDataFailed.Wrap(err, "file to record data")
		return models.Comment{}, err
	}

	return newComment, nil
}
