async function deleteComment(elementRef) {

	const commentID = elementRef.dataset.commentId
	if (!confirm("Are you sure you want to delete comment?")) {
		return; 
	}
	let response; 
	try {
		response = await fetch('/comments/' + commentID, { method: 'DELETE' });
	} catch(err) {
		console.log('Error:', err); 
		return; 
	}

	if (response.ok) {
		location.reload();
        return;
	}
    handleDeleteError(response);
}

async function deleteThread(elementRef) {

    const categoryID = elementRef.dataset.categoryId
    const threadID = elementRef.dataset.threadId
    if (!confirm("Confirm to delete thread")) {
        return;
    }

    let response;
    try {
        response = await fetch('/threads/' + threadID, { method: 'DELETE' }); 
    } catch(err) {
        console.log('Error:', err);
        return; 
    }

    if (response.ok) {
         location.href = '/categories/' + categoryID; // go to category page of where the deleted thread was from
         return;
    } 
    handleDeleteError(response);
}

function handleDeleteError(response) {
    if (response.status === 400) {
		window.alert("Invalid request");
	} else if (response.status === 401) {
		location.href = '/users/login';
	} else if (response.status === 404) {
		window.alert("Resource not found");
	} else {
		window.alert("Something went wrong. Please try again later.");
	}
}