const tagButtons = document.querySelectorAll('.tag-btn');

for (const tagButton of tagButtons) {
    tagButton.addEventListener('click', () => {
        if (tagButton.dataset.selected === 'true') {
            tagButton.dataset.selected = 'false';
        } else {
            tagButton.dataset.selected = 'true';
        }
    });
}
