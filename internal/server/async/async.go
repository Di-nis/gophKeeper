package async


import (
	"context"
	"github.com/Di-nis/gophKeeper/internal/model"
)

// TODO: как сделать?
// Publusher - публикатор.
type Publusher struct {

}

// generator - генерирует сообщения в канал.
func (urlUseCase *URLUseCase) generator(ctx context.Context, urls []models.URLBase, inChan chan models.URLBase) {
	for _, url := range urls {
		select {
		case <-ctx.Done():
			return
		case inChan <- url:
		}
	}
}


// worker - работник.
func (urlUseCase *URLUseCase) worker(ctx context.Context, urls <-chan models.URLBase, result chan error) {
	urlsToDB := make([]models.URLBase, 0, 100)

	for {
		select {
		case <-ctx.Done():
			if len(urlsToDB) > 0 {
				result <- urlUseCase.Repo.Delete(ctx, urlsToDB)
			}
			return

		case url, ok := <-urls:
			if !ok {
				if len(urlsToDB) > 0 {
					result <- urlUseCase.Repo.Delete(ctx, urlsToDB)
				}
				return
			}
			urlsToDB = append(urlsToDB, url)
			if len(urlsToDB) >= 1 {
				result <- urlUseCase.Repo.Delete(ctx, urlsToDB)
				urlsToDB = urlsToDB[:0]
			}

		}
	}

}