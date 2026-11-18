let items = [];
let editingItem = {};
let selectedCategory = '';

const TotalCostColumn = 6;
const CategoryColumn = 10;

const tableKeys = ["number", "quantity", "name", "description", "value", "total_cost",
    "manufacturer", "model", "category", "purchase_date", "serial_number"];
const itemDetailsKeys = ["quantity", "name", "description", "value", "total_cost", "manufacturer",
    "model", "category", "purchase_date", "serial_number", "details"];
const dollarColumns = ["value", "total_cost"];

// TODO: add total number of items near total cost

function renderTable(items) {
    const tHead = document.getElementById('itemsTableHead');
    const headRow = document.createElement('tr');

    items.sort((a, b) => {
        return a.name.localeCompare(b.name) ||
            a.description.localeCompare(b.description);
    });

    tHead.replaceChildren();
    tableKeys.forEach(key => {
        let el = document.createElement('th');
        el.textContent = getItemLabel(key);
        headRow.appendChild(el);
    });
    tHead.appendChild(headRow);

    const tbody = document.getElementById('itemsTableBody');
    tbody.replaceChildren();

    for (let i = 0; i < items.length; i++) {
        const row = document.createElement('tr');
        let item = items[i];

        tableKeys.forEach(key => {
            let el = document.createElement('td');
            if (dollarColumns.includes(key)) {
                el.append(accountingCell(item[key]));
            } else if (key === 'number') {
                el.textContent = String(i + 1);
            } else {
                el.textContent = item[key];
            }
            row.appendChild(el);
        });

        row.addEventListener('click', () => {
            showItemDetail(item, true);
        });
        tbody.appendChild(row);
    }
}

function accountingCell(amount) {
    const td = document.createElement('td');
    td.className = 'accounting';
    let el = document.createElement('span');
    el.className = "currency";
    el.textContent = '$';
    td.appendChild(el);

    el = document.createElement('span');
    el.className = "amount";
    el.textContent = amount;
    td.appendChild(el);
    return td;
}

function getItemIndexByID(id) {
    return items.findIndex(item => item.id === id);
}

function showItemDetail(item, disabled) {
    editingItem = {...item};
    itemDetailsKeys.forEach(key => {
        const el = document.getElementById(getItemKeyID(key));
        if (el) {
            el.value = item[key];
            el.disabled = disabled;
        }
    });
    viewButtons()
    document.getElementById('itemModal').classList.remove('hidden');
}

function closeItemDetails() {

    document.getElementById('itemModal').classList.add('hidden');
}

function editItemDetails() {
    showItemDetail(editingItem, false);
    editButtons();
}

function saveItemDetails() {
    let saveItem = {};

    itemDetailsKeys.forEach(key => {
        const el = document.getElementById(getItemKeyID(key));
        if (el) {
            switch (key) {
                case 'quantity':
                case 'value':
                    saveItem[key] = parseInt(el.value, 10);
                    console.log(saveItem.value);
                    console.log(typeof saveItem.value);
                    editingItem[key] = el.value;
                    break;
                case 'purchase_date':
                    if (!el.value) {
                        saveItem[key] = null;
                        editingItem[key] = null;
                        break;
                    }
                    const [month, day, year] = el.value.split("/").map(Number);
                    const date = new Date(Date.UTC(year, month - 1, day));
                    saveItem[key] = date.toISOString();
                    break;
                default:
                    saveItem[key] = el.value;
                    editingItem[key] = el.value;
            }
        }
    });
    saveItem['id'] = editingItem['id'];
    delete saveItem['number'];
    delete saveItem['total_cost'];

    fetch('/api/save', {
        method: 'POST',
        headers: {'Content-Type': 'application/json'},
        body: JSON.stringify(saveItem)
    })
        .then(response => {
            if (!response.ok) {
                throw new Error('Save failed: ' + response.status);
            }
            let index = getItemIndexByID(editingItem.id);
            items[index] = {...editingItem};
            closeItemDetails();
            showSelectedCategory(selectedCategory);
        })
        .catch(err => {
            console.error(err);
            alert('Save failed');
        });
}

function getItemKeyID(key) {
    return 'item' + getItemLabel(key)
        .replace(/\s/g, '');
}

function getItemLabel(key) {
    switch (key) {
        case 'quantity':
            return 'Qty';
        case 'number':
            return 'Num';
        case 'value_dollars':
            return 'Value';
        case 'total_cost_dollars':
            return 'Total Cost';
        default:
            return key
                .replace(/_/g, ' ')
                .replace(/\b\w/g, char => char.toUpperCase());
    }
}

function viewButtons() {
    document.getElementById('modalButtonClose').classList.remove('hidden');
    document.getElementById('modalButtonEdit').classList.remove('hidden');
    document.getElementById('modalButtonCancel').classList.add('hidden');
    document.getElementById('modalButtonSave').classList.add('hidden');
}

function editButtons() {
    document.getElementById('modalButtonClose').classList.add('hidden');
    document.getElementById('modalButtonEdit').classList.add('hidden');
    document.getElementById('modalButtonCancel').classList.remove('hidden');
    document.getElementById('modalButtonSave').classList.remove('hidden');
}

function showSelectedCategory(category) {
    selectedCategory = category;
    if (selectedCategory === '') {
        renderTable(items);
        return;
    }
    let displayItems = [];
    let totalCost = 0;

    items.forEach(function (item) {
        if (item.category === '' || item.category === selectedCategory) {
            totalCost += item.total_cost;
            displayItems.push(item);
        }
    });
    renderTable(displayItems);
    document.getElementById('totalAmount').textContent = "$ " + totalCost.toLocaleString('en-US');
}

function createDetailsTable() {
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
}

fetch('/api/items')
    .then(response => response.json())
    .then(data => {
        items = data;
        renderTable(items);
    });

document.addEventListener('DOMContentLoaded', () => {
    document.getElementById('categoryFilter').addEventListener('change', (event) => {
        showSelectedCategory(event.target.value)
    });
    createDetailsTable();
});

