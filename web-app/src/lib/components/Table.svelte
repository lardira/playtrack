<script lang="ts">
    export let rows;
</script>

<table>
    <thead>
        <tr>
            <th>Место</th>
            <th>Пользователь</th>
            <th>Текущая игра</th>
            <th>Поинты</th>
            <th>Пройдено</th>
            <th>Дроп</th>
            <th>Реролл</th>
        </tr>
    </thead>

    <tbody>
        {#each rows as row}
            <tr>
                <td>{row.place}</td>
                <td>
                    <a href={`/users/${row.user_id}`}>{row.nickname}</a>
                </td>
                <td>{row.current_game}</td>
                <td>{row.points}</td>
                <td>{row.completed}</td>
                <td>{row.dropped}</td>
                <td>{row.reroll}</td>
            </tr>
        {/each}
    </tbody>
</table>
