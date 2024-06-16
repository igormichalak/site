const tagButtons = document.querySelectorAll('.tag-btn');
const searchInput = document.querySelector('input#search');

const parser = new DOMParser();

function search() {
    const query = searchInput.value;
    const tags = [...tagButtons]
        .filter(el => el.dataset.selected === 'true')
        .map(el => el.dataset.tag);

    const url = new URL('/search', window.location.origin);

    url.searchParams.set('q', query);

    for (const tag of tags) {
        url.searchParams.append('tags', tag);
    }

    fetch(url.href).then(res => res.text()).then(htmlString => {
        const doc = parser.parseFromString(htmlString, 'text/html');
        document.querySelector('.post-list').replaceWith(doc.body.firstChild);
    });
}

for (const tagButton of tagButtons) {
    tagButton.addEventListener('click', () => {
        if (tagButton.dataset.selected === 'true') {
            tagButton.dataset.selected = 'false';
        } else {
            tagButton.dataset.selected = 'true';
        }
        search();
    });
}

searchInput.addEventListener('input', () => {
    search();
});
