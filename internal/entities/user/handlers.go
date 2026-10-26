package user

import (
	"forum/internal/errs"
	"forum/internal/pagination"
	"forum/internal/password"
	"forum/internal/render"
	"forum/internal/session"
	"net/http"
	"strings"
)

type UserHandler struct {
	service  *UserService
	sm       *session.SessionManager
	renderer *render.Renderer
}

func NewHandler(service *UserService, sm *session.SessionManager, renderer *render.Renderer) *UserHandler {
	return &UserHandler{
		service:  service,
		sm:       sm,
		renderer: renderer,
	}
}

func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user := GetUser(r)
	if user == nil {
		http.Redirect(w, r, "/users/login", http.StatusSeeOther)
		return
	}

	profileUsername := r.PathValue("username")

	userProfile, err := h.service.GetUserProfileByUsername(ctx, profileUsername)
	if err != nil {
		h.renderer.WriteErrorPage(w, err, render.ErrorPage{Username: user.Username})
		return
	}

	data := ProfilePage{
		ProfileUsername: profileUsername,
		Username:        user.Username,
		Role:            "Member",
		DateCreated:     userProfile.CreatedAt, //userProfile.CreatedAt.Format("02 Jan 2006"),
	}

	h.renderer.RenderPage(w, "profile.html", data)
}

func (h *UserHandler) RegisterPage(w http.ResponseWriter, r *http.Request) {
	// Redirect to home page if user is already logged in
	session := session.GetSession(r)
	if session != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	h.renderer.RenderPage(w, "register.html", RegistrationPage{})
}

func (h *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if err := r.ParseForm(); err != nil {
		h.renderer.WriteErrorPage(w, errs.ErrInvalidUserInput, render.ErrorPage{})
		return
	}

	inputUsername := strings.TrimSpace(r.PostFormValue("username"))
	inputEmail := strings.TrimSpace(r.PostFormValue("email"))

	user := &User{
		Username: inputUsername,
		Email:    inputEmail,
		Password: password.Password{},
		RoleID:   1,
	}

	user, err := h.service.RegisterUser(ctx, user, r.PostFormValue("password"))
	if err != nil {
		pageData := RegistrationPage{InputUsername: inputUsername, InputEmail: inputEmail} // service.RegisterUser returns a nil user on error, so values are taken from the form
		fieldErrors, statusCode := getRegistrationFieldErrors(err)
		pageData.FieldErrors = fieldErrors
		w.WriteHeader(statusCode)
		h.renderer.RenderPage(w, "register.html", pageData)
		return
	}

	// If registration successful, login user (create new session) and redirect to home
	if err := h.sm.Login(user.ID, w, r); err != nil {
		h.renderer.WriteErrorPage(w, err, render.ErrorPage{Username: user.Username})
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *UserHandler) LoginPage(w http.ResponseWriter, r *http.Request) {
	// Redirect to home page if user is already logged in
	session := session.GetSession(r)
	if session != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	h.renderer.RenderPage(w, "login.html", LoginPage{})
}
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Redirect to home page if user is already logged in
	session := session.GetSession(r)
	if session != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if err := r.ParseForm(); err != nil {
		h.renderer.WriteErrorPage(w, err, render.ErrorPage{})
		return
	}

	inputEmail := strings.TrimSpace(r.PostFormValue("email"))
	inputPassword := r.PostFormValue("password")

	input := LoginSubmission{
		Email:    inputEmail,
		Password: inputPassword,
	}

	userID, err := h.service.Authenticate(ctx, input)
	if err != nil {
		pageData := LoginPage{Email: inputEmail}
		fieldErrors, statusCode := getLoginFieldErrors(err)
		pageData.FieldErrors = fieldErrors
		w.WriteHeader(statusCode)
		h.renderer.RenderPage(w, "login.html", pageData)
		return
	}

	// Create new session
	if err := h.sm.Login(userID, w, r); err != nil {
		h.renderer.WriteErrorPage(w, err, render.ErrorPage{})
		return
	}

	// Redirect to home on successful login
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if err := h.sm.Logout(w, r); err != nil {
		h.renderer.WriteErrorPage(w, err, render.ErrorPage{})
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// ProfileThreadsPage gets all threads posted by a user and renders the data on the user's profile page.
func (h *UserHandler) ProfileThreadsPage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user := GetUser(r)
	if user == nil {
		http.Redirect(w, r, "/users/login", http.StatusSeeOther)
		return
	}

	type ThreadsPage struct {
		ProfileUsername string
		Username        string
		Threads         []ProfileThread
		PageNav         pagination.PaginationNav
	}

	paginationData, err := pagination.Parse(r.URL.Query())
	if err != nil {
		h.renderer.WriteErrorPage(w, errs.ErrInvalidUserInput, render.ErrorPage{Username: user.Username})
		return
	}

	if paginationData.Validate() != nil {
		h.renderer.WriteErrorPage(w, errs.ErrInvalidUserInput, render.ErrorPage{Username: user.Username})
		return
	}

	profileUsername := r.PathValue("username")

	// Get all the threads posted by a user
	threads, totalCount, err := h.service.GetThreadsByUsername(ctx, profileUsername, paginationData)
	if err != nil {
		h.renderer.WriteErrorPage(w, err, render.ErrorPage{Username: user.Username})
		return
	}

	totalPages := pagination.CalcTotalPages(totalCount, paginationData.PageSize)
	pageList := pagination.GetDisplayPageNums(paginationData.Page, totalPages, 5) // The visible pages in navigation

	pageNav := pagination.PaginationNav{
		Path:       r.URL.Path,
		Page:       paginationData.Page,
		TotalPages: totalPages,
		PageList:   pageList,
		PrevPage:   paginationData.Page - 1,
		NextPage:   paginationData.Page + 1,
	}
	pageData := ThreadsPage{
		ProfileUsername: profileUsername,
		Username:        user.Username, // logged in user
		Threads:         threads,
		PageNav:         pageNav,
	}

	h.renderer.RenderPage(w, "profile-threads.html", pageData)
}

func (h *UserHandler) ProfileCommentsPage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user := GetUser(r)
	if user == nil {
		http.Redirect(w, r, "/users/login", http.StatusSeeOther)
		return
	}

	type CommentsPage struct {
		ProfileUsername string
		Username        string
		Comments        []ProfileComment
		PageNav         pagination.PaginationNav
	}

	paginationData, err := pagination.Parse(r.URL.Query())
	if err != nil {
		h.renderer.WriteErrorPage(w, errs.ErrInvalidUserInput, render.ErrorPage{Username: user.Username})
		return
	}

	if paginationData.Validate() != nil {
		h.renderer.WriteErrorPage(w, errs.ErrInvalidUserInput, render.ErrorPage{Username: user.Username})
		return
	}

	profileUsername := r.PathValue("username")

	//Get all comments posted by a user
	comments, totalCount, err := h.service.GetCommentsByUsername(ctx, profileUsername, paginationData)
	if err != nil {
		h.renderer.WriteErrorPage(w, err, render.ErrorPage{Username: user.Username})
		return
	}

	totalPages := pagination.CalcTotalPages(totalCount, paginationData.PageSize)
	pageList := pagination.GetDisplayPageNums(paginationData.Page, totalPages, 5) //5 visible pages in navigation

	pageNav := pagination.PaginationNav{
		Path:       r.URL.Path,
		Page:       paginationData.Page,
		TotalPages: totalPages,
		PageList:   pageList,
		PrevPage:   paginationData.Page - 1,
		NextPage:   paginationData.Page + 1,
	}
	pageData := CommentsPage{
		ProfileUsername: profileUsername,
		Username:        user.Username,
		Comments:        comments,
		PageNav:         pageNav,
	}

	h.renderer.RenderPage(w, "profile-comments.html", pageData)

}

// ProfileLikedThreadsPage gets all liked threads by a user and renders the data on the user's profile page.
func (h *UserHandler) ProfileLikedThreadsPage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user := GetUser(r)
	if user == nil {
		http.Redirect(w, r, "/users/login", http.StatusSeeOther)
		return
	}

	type LikesPage struct {
		ProfileUsername string
		Username        string
		Threads         []ProfileThread
		PageNav         pagination.PaginationNav
	}

	paginationData, err := pagination.Parse(r.URL.Query())
	if err != nil {
		h.renderer.WriteErrorPage(w, errs.ErrInvalidUserInput, render.ErrorPage{Username: user.Username})
		return
	}

	if paginationData.Validate() != nil {
		h.renderer.WriteErrorPage(w, errs.ErrInvalidUserInput, render.ErrorPage{Username: user.Username})
		return
	}

	profileUsername := r.PathValue("username")

	// Get all the liked threads posted by a user
	threads, totalCount, err := h.service.GetLikedThreadsByUsername(ctx, profileUsername, paginationData)
	if err != nil {
		h.renderer.WriteErrorPage(w, err, render.ErrorPage{Username: user.Username})
		return
	}

	totalPages := pagination.CalcTotalPages(totalCount, paginationData.PageSize)
	pageList := pagination.GetDisplayPageNums(paginationData.Page, totalPages, 5) // The visible pages in navigation

	pageNav := pagination.PaginationNav{
		Path:       r.URL.Path,
		Page:       paginationData.Page,
		TotalPages: totalPages,
		PageList:   pageList,
		PrevPage:   paginationData.Page - 1,
		NextPage:   paginationData.Page + 1,
	}

	pageData := LikesPage{
		ProfileUsername: profileUsername,
		Username:        user.Username,
		Threads:         threads,
		PageNav:         pageNav,
	}

	h.renderer.RenderPage(w, "profile-liked-threads.html", pageData)
}

// ProfileDislikedThreadsPage gets all disliked threads by a user and renders the data on the user's profile page.
func (h *UserHandler) ProfileDislikedThreadsPage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user := GetUser(r)
	if user == nil {
		http.Redirect(w, r, "/users/login", http.StatusSeeOther)
		return
	}

	type DislikesPage struct {
		ProfileUsername string
		Username        string
		Threads         []ProfileThread
		PageNav         pagination.PaginationNav
	}

	paginationData, err := pagination.Parse(r.URL.Query())
	if err != nil {
		h.renderer.WriteErrorPage(w, errs.ErrInvalidUserInput, render.ErrorPage{Username: user.Username})
		return
	}

	if paginationData.Validate() != nil {
		h.renderer.WriteErrorPage(w, errs.ErrInvalidUserInput, render.ErrorPage{Username: user.Username})
		return
	}

	profileUsername := r.PathValue("username")

	// Get all the disliked threads posted by a user
	threads, totalCount, err := h.service.GetDislikedThreadsByUserID(ctx, profileUsername, paginationData)
	if err != nil {
		h.renderer.WriteErrorPage(w, err, render.ErrorPage{Username: user.Username})
		return
	}

	totalPages := pagination.CalcTotalPages(totalCount, paginationData.PageSize)
	pageList := pagination.GetDisplayPageNums(paginationData.Page, totalPages, 5) // The visible pages in navigation

	pageNav := pagination.PaginationNav{
		Path:       r.URL.Path,
		Page:       paginationData.Page,
		TotalPages: totalPages,
		PageList:   pageList,
		PrevPage:   paginationData.Page - 1,
		NextPage:   paginationData.Page + 1,
	}

	pageData := DislikesPage{
		ProfileUsername: profileUsername,
		Username:        user.Username,
		Threads:         threads,
		PageNav:         pageNav,
	}

	h.renderer.RenderPage(w, "profile-disliked-threads.html", pageData)
}
