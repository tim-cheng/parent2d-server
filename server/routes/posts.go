package routes

import (
	"github.com/codegangsta/martini"
	"github.com/codegangsta/martini-contrib/render"
	"github.com/tim-cheng/mila/server/models"
	"io/ioutil"
	"net/http"
	"strconv"
)

func (rt *Routes) PostPost(r render.Render, req *http.Request) {
	post, err := rt.Db.NewPost(req.FormValue("user_id"), req.FormValue("body"))
	if err == nil {
		err = rt.Db.PostPost(post)
	}
	if err == nil {
		r.JSON(201, map[string]interface{}{
			"id": post.Id,
		})
	} else {
		r.JSON(404, map[string]interface{}{
			"message": "Failed to add post " + err.Error(),
		})
	}
}

func (rt *Routes) GetPosts(r render.Render, req *http.Request) {
	posts, err := rt.Db.GetPosts(req.FormValue("user_id"), req.FormValue("degree"))
	if err == nil && len(posts) > 0 {
		retPosts := make([]map[string]interface{}, len(posts))
		for i := range posts {
			p := posts[i].(*models.Post)
			nComments, _ := rt.Db.GetNumComments(p.Id)
			nStars, _ := rt.Db.GetNumStars(p.Id)
			nSelfStar, _ := rt.Db.GetStarByUser(p.Id, req.FormValue("user_id"))
			retPosts[i] = map[string]interface{}{
				"id":           p.Id,
				"user_id":      p.UserId,
				"body":         p.Body,
				"created_at":   p.CreatedAt,
				"num_comments": nComments,
				"num_stars":    nStars,
				"self_star":    nSelfStar,
				"has_picture":  p.HasPicture,
			}
		}
		r.JSON(200, retPosts)
	} else {
		r.JSON(404, map[string]interface{}{
			"message": "Failed to get posts ",
		})
	}
}

func (rt *Routes) PostPostPicture(params martini.Params, r render.Render, req *http.Request) {
	buf, err := ioutil.ReadAll(req.Body)
	if err != nil {
		r.JSON(404, map[string]interface{}{
			"message": "failed to read picture " + err.Error(),
		})
		return
	}
	err = rt.Db.PostPostPicture(params["id"], buf)
	if err == nil {
		r.JSON(201, map[string]interface{}{"id": params["id"]})
	} else {
		r.JSON(404, map[string]interface{}{
			"message": "failed to save post picture " + err.Error(),
		})
	}
}

func (rt *Routes) GetPostPicture(params martini.Params, r render.Render, w http.ResponseWriter) {
	image, err := rt.Db.GetPostPicture(params["id"])
	if err == nil {
		w.Header().Set("Content-Type", "image/jpeg")
		w.Header().Set("Content-Length", strconv.Itoa(len(image)))
		w.Write(image)
	} else {
		r.JSON(404, map[string]interface{}{
			"message": "failed to retrieve post picture " + err.Error(),
		})
	}
}
