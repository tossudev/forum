// Editor's note:
// This code is really spaghetti-y and could be fixed in the future
// also worth considering just fetching the data from backend after action
// although it's better UX to update the likes realtime

// Element name in HTML is comment ID
// use it to send like post req to backend
async function likeComment(elementRef, like) {
	const commentName = elementRef.name
	const commentID = elementRef.dataset.commentId

	const likeButton = document.querySelector(
	  `button[data-comment-id="${commentID}"][data-action="like"]`
	);

	const dislikeButton = document.querySelector(
	  `button[data-comment-id="${commentID}"][data-action="dislike"]`
	);

	const likeCounter = document.querySelector(`span.likes-${commentID}`)
	const dislikeCounter = document.querySelector(`span.dislikes-${commentID}`)

	let initialLikes = Number(likeCounter.textContent)
	let initialDislikes = Number(dislikeCounter.textContent)

	if (likeButton.classList.contains("reaction-btn-pressed")) {
		initialLikes -= 1
	}
	if (dislikeButton.classList.contains("reaction-btn-pressed")) {
		initialDislikes -= 1
	}

    let response;
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
		likeCounter.textContent = Number(initialLikes)
		dislikeCounter.textContent = Number(initialDislikes)

		likeCounter.name = "false"
		dislikeCounter.name = "false"

		let addReaction = elementRef.classList.contains("reaction-btn")

		likeButton.classList.remove("reaction-btn-pressed")
		dislikeButton.classList.remove("reaction-btn-pressed")

		likeButton.classList.add("reaction-btn")
		dislikeButton.classList.add("reaction-btn")

		if (addReaction) {
			if (like) {
				likeCounter.textContent = Number(initialLikes) + 1
				likeCounter.name = "true"
			} else {
				dislikeCounter.textContent = Number(initialDislikes) + 1
				dislikeCounter.name = "true"
			}
			elementRef.classList.remove("reaction-btn")
			elementRef.classList.add("reaction-btn-pressed")
		}
    }
}
