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
	}
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
        location.href = '/categories/' + categoryID;
    } else if (response.status === 401) {
        location.href = '/users/login';
    } else {
        window.alert("Something went wrong. Please try again later.");
    }
}
