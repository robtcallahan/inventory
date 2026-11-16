fetch('/api/items')
    .then(response => response.json())
    .then(items => {
        const tbody = document.getElementById('items');
        items.forEach(item => {
            const row = document.createElement('tr');
            const cell = document.createElement('td');
            cell.textContent = item.Name;
            row.appendChild(cell);
            tbody.appendChild(row);
        });
    });
