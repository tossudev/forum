// Editor's note:
// This code is really spaghetti-y and could be fixed in the future
// also worth considering just fetching the data from backend after action
// although it's better UX to update the likes realtime

// Element name in HTML is comment ID
// use it to send like post req to backend
async function likeComment(elementRef, like) {
    let response;
	var commentName = elementRef.name

	let matches = commentName.match(/(\d+)/);
	let commentID = matches[0]

	let likeButton = document.querySelector('button[name="likebtn-' + commentID + '"]')
	let dislikeButton = document.querySelector('button[name="dislikebtn-' + commentID + '"]')

	let likeCounter = document.querySelector('span.likes-' + commentID)
	let dislikeCounter = document.querySelector('span.dislikes-' + commentID)

	let likeAmount = Number(likeCounter.innerHTML)
	let dislikeAmount = Number(dislikeCounter.innerHTML)

	if (likeButton.classList.contains("reaction-btn-pressed")) {
		likeAmount -= 1
	}
	if (dislikeButton.classList.contains("reaction-btn-pressed")) {
		dislikeAmount -= 1
	}

    try {
        response = await fetch('/like/comment', {
            method: 'POST',
            body: JSON.stringify({
				"comment_id": Number(commentID),
				"like": like
			})
        });
    } catch(err) {
        console.log('Error:', err);
    }

    if (response.ok) {
		likeCounter.innerHTML = Number(likeAmount)
		dislikeCounter.innerHTML = Number(dislikeAmount)

		likeCounter.name = "false"
		dislikeCounter.name = "false"

		let addReaction = elementRef.classList.contains("reaction-btn")

		likeButton.classList.remove("reaction-btn-pressed")
		dislikeButton.classList.remove("reaction-btn-pressed")

		likeButton.classList.add("reaction-btn")
		dislikeButton.classList.add("reaction-btn")

		if (addReaction) {
			if (like) {
				likeCounter.innerHTML = Number(likeAmount) + 1
				likeCounter.name = "true"
			} else {
				dislikeCounter.innerHTML = Number(dislikeAmount) + 1
				dislikeCounter.name = "true"
			}
			elementRef.classList.remove("reaction-btn")
			elementRef.classList.add("reaction-btn-pressed")
		}
    }
}
