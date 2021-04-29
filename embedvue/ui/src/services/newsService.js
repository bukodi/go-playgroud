const url = '/data/articles.json';
class NewsService {
    getArticlesByType(articleType) {
        return fetch(url)
            .then((response) => {
            return response.json();
        })
            .then((serverArticles) => {
            const newsArticles = serverArticles
                .filter((serverArticle) => serverArticle.articleType === articleType)
                .map(NewsService.map);
            return newsArticles;
        })
            .catch((e) => {
            console.error('An error occurred retrieving the news articles from ' + url, e);
        });
    }
    getFavorites() {
        return fetch(url)
            .then((response) => {
            return response.json();
        })
            .then((serverArticles) => {
            const newsArticles = serverArticles
                .filter((serverArticle) => serverArticle.isFavourite === true)
                .map(NewsService.map);
            return newsArticles;
        })
            .catch((e) => {
            console.error('An error occurred retrieving the news articles from ' + url, e);
        });
    }
    static map(serverArticle) {
        return {
            id: serverArticle.id,
            title: serverArticle.title,
            content: serverArticle.content,
            dateString: serverArticle.dateString,
            baseImageName: serverArticle.baseImageName,
            articleType: serverArticle.articleType,
            isFavourite: serverArticle.isFavourite
        };
    }
}
export default new NewsService();
//# sourceMappingURL=newsService.js.map