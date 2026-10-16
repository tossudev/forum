// Send like request to server and wait for approval
// Upon success, update the frontend to match user action

// Frontend acts like a fake mirror,
// it does not request the real data from the server
// but since it is only updated on server acceptance,
// it should always display correctly

async function likeComment(elementRef, like) {
	const commentID = elementRef.dataset.commentId

	// Send request to server and await response
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

	// Action accepted by server, update frontend
    if (response.ok) {
		const likeButton = document.querySelector(
		  `button[data-comment-id="${commentID}"][data-action="like"]`
		);

		const dislikeButton = document.querySelector(
		  `button[data-comment-id="${commentID}"][data-action="dislike"]`
		);

		const likeCounter = document.querySelector(`span.likes-${commentID}`)
		const dislikeCounter = document.querySelector(`span.dislikes-${commentID}`)

		// Likes and dislikes without current users input
		// eg. initial values before action
		let initialLikes = Number(likeCounter.textContent)
		let initialDislikes = Number(dislikeCounter.textContent)
		if (likeButton.classList.contains("pressed")) {
			initialLikes -= 1
		}
		if (dislikeButton.classList.contains("pressed")) {
			initialDislikes -= 1
		}

		let addReaction = !elementRef.classList.contains("pressed")

		// Reset buttons & counters to intial states
		likeCounter.textContent = Number(initialLikes)
		dislikeCounter.textContent = Number(initialDislikes)

		likeButton.classList.remove("pressed")
		dislikeButton.classList.remove("pressed")

		// If user removed like/dislike we only need to reset buttons
		if (!addReaction) {
			return;
		}

		if (like) {
			likeCounter.textContent = Number(initialLikes) + 1
		} else {
			dislikeCounter.textContent = Number(initialDislikes) + 1
		}
		elementRef.classList.add("pressed")
    }
}
