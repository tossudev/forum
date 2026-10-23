// Send like request to server and wait for approval
// Upon success, update the frontend to match user action

// Frontend acts like a fake mirror,
// it does not request the real data from the server
// but since it is only updated on server acceptance,
// it should always display correctly

async function like(elementRef, like) {
	// Assign id and POST path based on if the like was on a comment vs thread
	const isComment = !!elementRef.dataset.commentId;
	const id = isComment ? elementRef.dataset.commentId : elementRef.dataset.threadId;
	const path = isComment ? '/like/comment' : '/like/thread';
	const dataName = isComment ? "comment" : "thread";

	// Send request to server and await response
    let response;
    try {
        response = await fetch(path, {
            method: 'POST',
            body: JSON.stringify({
				"id": Number(id),
				"like": like
			})
        });
    } catch(err) {
		window.alert("Server Error: couldn't send like");
        console.log('Error:', err);
    }

	// Action accepted by server, update frontend
    if (response.ok) {
		const likeButton = document.querySelector(
		  `button[data-${dataName}-id="${id}"][data-action="like"]`
		);

		const dislikeButton = document.querySelector(
		  `button[data-${dataName}-id="${id}"][data-action="dislike"]`
		);

		const likeCounter = likeButton.querySelector(`span.like`)
		const dislikeCounter = dislikeButton.querySelector(`span.dislike`)

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
    } else {
		window.alert("You need to be signed in!");
	}
}
