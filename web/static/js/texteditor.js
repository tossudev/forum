function toggleEditorVisible() {
    const editor = document.querySelector('form.texteditor')
    const listing = document.querySelector('table.listing-container')
    const createButton = document.querySelector('button.create-btn')

	const showEditor = editor.hidden;

	editor.hidden = !showEditor;
	listing.hidden = showEditor;
	createButton.hidden = showEditor;
}
