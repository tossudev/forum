
// Element name in HTML is comment ID
// use it to send like post req to backend
async function likeComment(elementRef, like) {
    let response;

    try {
        response = await fetch('/like/comment', {
            method: 'POST',
            body: JSON.stringify({
				"comment_id": Number(elementRef.name),
				"like": like
			})
        });
    } catch(err) {
        console.log('Error:', err);
    }

    if (response.ok) {
		elementRef.innerHTML += 1;
        return;
    }
}

async function getLike(elementRef) {
    let response;

    try {
        response = await fetch('/like/comment', {
            method: 'GET',
            body: JSON.stringify({
				"comment_id": Number(elementRef.name),
				"like": like
			})
        });
    } catch(err) {
        console.log('Error:', err);
    }

    if (response.ok) {
		elementRef.innerHTML += 1;
        return;
    }

}
