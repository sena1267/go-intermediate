package services

import (
	"database/sql"
	"errors"
	"sync"

	"github.com/sena1267/go-intermediate/apperrors"
	"github.com/sena1267/go-intermediate/models"
	"github.com/sena1267/go-intermediate/repositories"
)

// PostArticleHandlerで使うことを想定したサービス
// 引数の情報をもとに新しい記事を作り、結果を返却
func (s *MyAppService) PostArticleService(article models.Article) (models.Article, error) {
	newArticle, err := repositories.InsertArticle(s.db, article)
	if err != nil {
		err = apperrors.InsertDataFailed.Wrap(err, "failed to record data")
		return models.Article{}, err

	}
	return newArticle, nil
}

// ArticleListHandlerで使うことを想定したサービス
// 指定pageの記事一覧を返却
func (s *MyAppService) GetArticleListService(page int) ([]models.Article, error) {
	articleList, err := repositories.SelectArticleList(s.db, page)
	if err != nil {
		err = apperrors.GetDataFailed.Wrap(err, "fail to get data")
		return nil, err
	}

	if len(articleList) == 0 {
		err = apperrors.NAData.Wrap(ErrNoData, "no data")
		return nil, err
	}

	return articleList, nil
}

// ArticleDetailHandlerで使うことを想定したサービス
// 指定IDの記事情報を返却
func (s *MyAppService) GetArticleService(articleID int) (models.Article, error) {
	var (
		article                      models.Article
		commentList                  []models.Comment
		articleGetErr, commentGeterr error

		amu sync.Mutex
		cmu sync.Mutex
		wg  sync.WaitGroup
	)

	wg.Add(2)

	go func(db *sql.DB, articleID int) {
		defer wg.Done()
		amu.Lock()
		article, articleGetErr = repositories.SelectArticleDetail(db, articleID)
		amu.Unlock()
	}(s.db, articleID)

	go func(db *sql.DB, articleID int) {
		defer wg.Done()
		cmu.Lock()
		commentList, commentGeterr = repositories.SelectCommentList(db, articleID)
		cmu.Unlock()
	}(s.db, articleID)

	wg.Wait()

	if articleGetErr != nil {
		if errors.Is(articleGetErr, sql.ErrNoRows) {
			articleGetErr = apperrors.NAData.Wrap(articleGetErr, "no data")
			return models.Article{}, articleGetErr
		}
		err := apperrors.GetDataFailed.Wrap(articleGetErr, "file to get data")
		return models.Article{}, err
	}

	if commentGeterr != nil {
		err := apperrors.GetDataFailed.Wrap(commentGeterr, "fail to get data")
		return models.Article{}, err
	}

	article.CommentList = append(article.CommentList, commentList...)

	return article, nil
}

// PostNiceHandlerで使うことを想定したサービス
// 指定IDの記事のいいね数を+1して、結果を返却
func (s *MyAppService) PostNiceService(article models.Article) (models.Article, error) {
	err := repositories.UpdateNiceNum(s.db, article.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = apperrors.NoTargetData.Wrap(err, "does not exist target article")
			return models.Article{}, err
		}
		err = apperrors.UpdateDataFailed.Wrap(err, "fail to update nice count")
		return models.Article{}, err
	}

	return models.Article{
		ID:        article.ID,
		Title:     article.Title,
		Contents:  article.Contents,
		UserName:  article.UserName,
		NiceNum:   article.NiceNum + 1,
		CreatedAt: article.CreatedAt,
	}, nil
}
