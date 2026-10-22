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
