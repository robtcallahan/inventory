fetch('/api/items')
    .then(response => response.json())
    .then(items => {
        const tHead = document.getElementById('itemsTableHead');

        const headRow = document.createElement('tr');

        const qtyHead = document.createElement('th');
        qtyHead.textContent = "Qty";
        headRow.appendChild(qtyHead);

        const nameHead = document.createElement('th');
        nameHead.textContent = "Name";
        headRow.appendChild(nameHead);

        const descrHead = document.createElement('th');
        descrHead.textContent = "Description";
        headRow.appendChild(descrHead);

        const valueHead = document.createElement('th');
        valueHead.textContent = "Value";
        headRow.appendChild(valueHead);

        const totCostHead = document.createElement('th');
        totCostHead.textContent = "Total Cost";
        headRow.appendChild(totCostHead);

        const manHead = document.createElement('th');
        manHead.textContent = "Manufacturer";
        headRow.appendChild(manHead);

        const modelHead = document.createElement('th');
        modelHead.textContent = "Model";
        headRow.appendChild(modelHead);

        const catHead = document.createElement('th');
        catHead.textContent = "Category";
        headRow.appendChild(catHead);

        const pDateHead = document.createElement('th');
        pDateHead.textContent = "Purchase Date";
        headRow.appendChild(pDateHead);

        const serNumHead = document.createElement('th');
        serNumHead.textContent = "Serial Number";
        headRow.appendChild(serNumHead);

        tHead.appendChild(headRow);

        const tbody = document.getElementById('itemsTableBody');
        items.forEach(item => {
            const row = document.createElement('tr');

            const qtyCell = document.createElement('td');
            qtyCell.textContent = item.quantity;
            row.appendChild(qtyCell);

            const nameCell = document.createElement('td');
            nameCell.textContent = item.name;
            row.appendChild(nameCell);

            const descrCell = document.createElement('td');
            descrCell.textContent = item.description;
            row.appendChild(descrCell);

            const valueCell = document.createElement('td');
            valueCell.textContent = item.value;
            valueCell.className = 'dollars';
            row.appendChild(valueCell);

            const totCostCell = document.createElement('td');
            totCostCell.textContent = item.total_cost;
            totCostCell.className = 'dollars';
            row.appendChild(totCostCell);

            const manCell = document.createElement('td');
            manCell.textContent = item.manufacturer;
            row.appendChild(manCell);

            const modelCell = document.createElement('td');
            modelCell.textContent = item.model;
            row.appendChild(modelCell);

            const catCell = document.createElement('td');
            catCell.textContent = item.category;
            row.appendChild(catCell);

            const pDateCell = document.createElement('td');
            pDateCell.textContent = item.purchase_date;
            row.appendChild(pDateCell);

            const serNumCell = document.createElement('td');
            serNumCell.textContent = item.serial_number;
            row.appendChild(serNumCell);

            tbody.appendChild(row);
        });
    });

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
});

