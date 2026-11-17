
let viewItem;
const tableKeys = ["number", "quantity", "name", "description", "value", "total_cost", "manufacturer",
    "model", "category", "purchase_date", "serial_number"];
const itemDetailsKeys = ["quantity", "name", "description", "value", "total_cost", "manufacturer",
    "model", "category", "purchase_date", "serial_number", "details"];

fetch('/api/items')
    .then(response => response.json())
    .then(items => {
        const tHead = document.getElementById('itemsTableHead');
        const headRow = document.createElement('tr');
        tableKeys.forEach(key => {
            let el = document.createElement('th');
            el.textContent = getItemLabel(key);
            headRow.appendChild(el);
        });
        tHead.appendChild(headRow);

        const tbody = document.getElementById('itemsTableBody');
        items.forEach(item => {
            const row = document.createElement('tr');

            tableKeys.forEach(key => {
                let el = document.createElement('td');
                el.textContent = item[key];
                if (key === 'value' || key === 'total_cost') {
                    el.className = 'dollars';
                }
                if (key === 'details') {
                    el.classList.add('hidden');
                }
                row.appendChild(el);
            });

            row.addEventListener('click', () => {
                showItemDetail(item);
            });
            tbody.appendChild(row);
        });
    });

function showItemDetail(item) {
    viewItem = item;
    itemDetailsKeys.forEach(key => {
        const el = document.getElementById(getItemKeyID(key));
        if (el) {
            el.value = item[key];
            el.disabled = true;
        }
    });

    document.getElementById('itemModal').classList.remove('hidden');
}

function closeItemDetails() {
    document.getElementById('itemModal').classList.add('hidden');
}

function editItemDetails() {
    itemDetailsKeys.forEach(key => {
        const el = document.getElementById(getItemKeyID(key));
        if (el) {
            // el.value = viewItem[key];
            el.disabled = false;
        }
    });

    document.getElementById('modalButtonClose').classList.add('hidden');
    document.getElementById('modalButtonEdit').classList.add('hidden');
    document.getElementById('modalButtonCancel').classList.remove('hidden');
    document.getElementById('modalButtonSave').classList.remove('hidden');
}

function saveItemDetails() {

}

function getItemKeyID(key) {
    return 'item' + getItemLabel(key)
        .replace(/\s/g, '');
}

function getItemLabel(key) {
    switch (key) {
        case 'quantity': return 'Qty';
        case 'number': return 'Num';
        default: return key
            .replace(/_/g, ' ')
            .replace(/\b\w/g, char => char.toUpperCase());
    }
}

document.addEventListener('DOMContentLoaded', () => {
    document.getElementById('categoryFilter').addEventListener('change', (event) => {
        const category = event.target.value;
        const tbody = document.getElementById('itemsTableBody');
        for (let i = 0; i < tbody.rows.length; i++) {
            const row = tbody.rows[i];
            const catCell = row.cells[7];
            if (category === '' || catCell.textContent === category) {
                row.style.display = '';
            } else {
                row.style.display = 'none';
            }
        }
    });

    const tbody = document.getElementById('itemDetailsBody');
    itemDetailsKeys.forEach(key => {
        const row = document.createElement('tr');
        let cell = document.createElement('td');
        const label = document.createElement('label');
        label.attributes['for'] = getItemKeyID(key);
        label.textContent = getItemLabel(key);
        cell.appendChild(label);
        row.appendChild(cell);

        cell = document.createElement('td');
        const input = document.createElement('input');
        input.id = getItemKeyID(key);
        input.type = 'text';
        cell.appendChild(input);
        row.appendChild(cell);
        tbody.appendChild(row);
    });
});

