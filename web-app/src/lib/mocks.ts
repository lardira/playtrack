import type { Player, LeaderboardPlayer, Game, PlayedGame } from './types';

export const playersMock: Player[] = [
    {
        id: '550e8400-e29b-41d4-a716-446655440000',
        username: 'Alice',
        img: null,
        email: 'alice@example.com',
        created_at: '2024-01-01T00:00:00Z'
    },
    {
        id: '0fe81b90-7e07-4c32-833d-25c3678afda9',
        username: 'Memeladon',
        img: null,
        email: 'memeladon@mail.ru',
        created_at: '2024-01-01T00:00:00Z'
    },
    {
        id: '550e8400-e29b-41d4-a716-446655440001',
        username: 'Bob',
        img: null,
        email: 'bob@example.com',
        created_at: '2024-01-02T00:00:00Z'
    },
    {
        id: '550e8400-e29b-41d4-a716-446655440002',
        username: 'Charlie',
        img: null,
        email: null,
        created_at: '2024-01-03T00:00:00Z'
    },
    {
        id: '550e8400-e29b-41d4-a716-446655440003',
        username: 'Diana',
        img: null,
        email: 'diana@example.com',
        created_at: '2024-01-04T00:00:00Z'
    },
    {
        id: '550e8400-e29b-41d4-a716-446655440004',
        username: 'Eve',
        img: '/favicon.png',
        email: 'eve@example.com',
        created_at: '2024-01-05T00:00:00Z'
    },
    {
        id: '550e8400-e29b-41d4-a716-446655440005',
        username: 'Frank',
        img: null,
        email: 'frank@example.com',
        created_at: '2024-01-06T00:00:00Z'
    },
    {
        id: '550e8400-e29b-41d4-a716-446655440006',
        username: 'Grace',
        img: null,
        email: null,
        created_at: '2024-01-07T00:00:00Z'
    },
    {
        id: '550e8400-e29b-41d4-a716-446655440007',
        username: 'Henry',
        img: null,
        email: 'henry@example.com',
        created_at: '2024-01-08T00:00:00Z'
    },
];

export const leaderboardMock: LeaderboardPlayer[] = [
    { player_id: '550e8400-e29b-41d4-a716-446655440000', total: 1200, completed: 15, dropped: 2, rerolled: 1 },
    { player_id: '550e8400-e29b-41d4-a716-446655440001', total: 900, completed: 10, dropped: 1, rerolled: 0 },
    { player_id: '550e8400-e29b-41d4-a716-446655440002', total: 850, completed: 12, dropped: 3, rerolled: 2 },
    { player_id: '550e8400-e29b-41d4-a716-446655440003', total: 750, completed: 8, dropped: 1, rerolled: 1 },
    { player_id: '550e8400-e29b-41d4-a716-446655440004', total: 680, completed: 9, dropped: 2, rerolled: 0 },
    { player_id: '550e8400-e29b-41d4-a716-446655440005', total: 620, completed: 7, dropped: 1, rerolled: 1 },
    { player_id: '550e8400-e29b-41d4-a716-446655440006', total: 580, completed: 6, dropped: 2, rerolled: 0 },
    { player_id: '550e8400-e29b-41d4-a716-446655440007', total: 520, completed: 5, dropped: 1, rerolled: 1 },
];

export const gamesMock: Game[] = [
    {
        id: 101,
        points: 3,
        hours_to_beat: 10,
        title: 'The Witcher 3',
        url: 'https://example.com/witcher3',
        created_at: '2024-01-01T00:00:00Z'
    },
    {
        id: 102,
        points: 5,
        hours_to_beat: 20,
        title: 'Dark Souls III',
        url: null,
        created_at: '2024-01-02T00:00:00Z'
    },
    {
        id: 103,
        points: 1,
        hours_to_beat: 2,
        title: 'Portal',
        url: 'https://example.com/portal',
        created_at: '2024-01-03T00:00:00Z'
    },
    {
        id: 104,
        points: 6,
        hours_to_beat: 24,
        title: 'Elden Ring',
        url: 'https://example.com/eldenring',
        created_at: '2024-01-04T00:00:00Z'
    },
    {
        id: 105,
        points: 4,
        hours_to_beat: 15,
        title: 'Hades',
        url: 'https://example.com/hades',
        created_at: '2024-01-05T00:00:00Z'
    },
    {
        id: 106,
        points: 8,
        hours_to_beat: 32,
        title: 'Factorio',
        url: null,
        created_at: '2024-01-06T00:00:00Z'
    },
    {
        id: 107,
        points: 2,
        hours_to_beat: 6,
        title: 'Celeste',
        url: 'https://example.com/celeste',
        created_at: '2024-01-07T00:00:00Z'
    },
    {
        id: 108,
        points: 7,
        hours_to_beat: 28,
        title: 'Baldur\'s Gate 3',
        url: 'https://example.com/bg3',
        created_at: '2024-01-08T00:00:00Z'
    },
    {
        id: 109,
        points: 3,
        hours_to_beat: 12,
        title: 'Hollow Knight',
        url: 'https://example.com/hollowknight',
        created_at: '2024-01-09T00:00:00Z'
    },
    {
        id: 110,
        points: 5,
        hours_to_beat: 18,
        title: 'Stardew Valley',
        url: null,
        created_at: '2024-01-10T00:00:00Z'
    },
];

export const gamesPlayedMock: PlayedGame[] = [
    {
        id: 1,
        player_id: '550e8400-e29b-41d4-a716-446655440000',
        game_id: 101,
        status: 'completed',
        points: 3,
        started_at: '2024-01-25T00:00:00Z',
        completed_at: '2024-01-26T00:00:00Z',
        comment: 'Отличная игра!',
        rating: 95,
        play_time: 'PT50H30M'
    },
    {
        id: 2,
        player_id: '550e8400-e29b-41d4-a716-446655440000',
        game_id: 102,
        status: 'in_progress',
        points: 5,
        started_at: '2024-01-27T00:00:00Z',
        completed_at: null,
        comment: null,
        rating: null,
        play_time: 'PT15H'
    },
    {
        id: 3,
        player_id: '550e8400-e29b-41d4-a716-446655440000',
        game_id: 103,
        status: 'completed',
        points: 1,
        started_at: '2024-01-28T00:00:00Z',
        completed_at: '2024-01-28T12:00:00Z',
        comment: 'Быстрая и интересная',
        rating: 88,
        play_time: 'PT2H15M'
    },
    {
        id: 4,
        player_id: '550e8400-e29b-41d4-a716-446655440001',
        game_id: 101,
        status: 'completed',
        points: 3,
        started_at: '2024-01-20T00:00:00Z',
        completed_at: '2024-01-22T00:00:00Z',
        comment: 'Лучшая RPG',
        rating: 98,
        play_time: 'PT45H'
    },
    {
        id: 5,
        player_id: '550e8400-e29b-41d4-a716-446655440001',
        game_id: 104,
        status: 'in_progress',
        points: 6,
        started_at: '2024-01-30T00:00:00Z',
        completed_at: null,
        comment: 'Сложная, но интересная',
        rating: null,
        play_time: 'PT30H'
    },
    {
        id: 6,
        player_id: '550e8400-e29b-41d4-a716-446655440002',
        game_id: 105,
        status: 'completed',
        points: 4,
        started_at: '2024-01-15T00:00:00Z',
        completed_at: '2024-01-18T00:00:00Z',
        comment: 'Отличный roguelike',
        rating: 92,
        play_time: 'PT18H45M'
    },
    {
        id: 7,
        player_id: '550e8400-e29b-41d4-a716-446655440002',
        game_id: 106,
        status: 'dropped',
        points: 8,
        started_at: '2024-01-10T00:00:00Z',
        completed_at: '2024-01-25T00:00:00Z',
        comment: 'Слишком сложно',
        rating: 65,
        play_time: 'PT25H'
    },
    {
        id: 8,
        player_id: '550e8400-e29b-41d4-a716-446655440003',
        game_id: 107,
        status: 'completed',
        points: 2,
        started_at: '2024-01-12T00:00:00Z',
        completed_at: '2024-01-13T00:00:00Z',
        comment: 'Прекрасная платформерка',
        rating: 90,
        play_time: 'PT8H20M'
    },
    {
        id: 9,
        player_id: '550e8400-e29b-41d4-a716-446655440003',
        game_id: 108,
        status: 'in_progress',
        points: 7,
        started_at: '2024-01-28T00:00:00Z',
        completed_at: null,
        comment: 'Очень долгая игра',
        rating: null,
        play_time: 'PT40H'
    },
    {
        id: 10,
        player_id: '550e8400-e29b-41d4-a716-446655440004',
        game_id: 109,
        status: 'completed',
        points: 3,
        started_at: '2024-01-05T00:00:00Z',
        completed_at: '2024-01-08T00:00:00Z',
        comment: 'Красивая метроидвания',
        rating: 94,
        play_time: 'PT15H30M'
    },
    {
        id: 11,
        player_id: '550e8400-e29b-41d4-a716-446655440004',
        game_id: 110,
        status: 'rerolled',
        points: 5,
        started_at: '2024-01-20T00:00:00Z',
        completed_at: '2024-01-22T00:00:00Z',
        comment: 'Надоело',
        rating: 70,
        play_time: 'PT12H'
    },
    {
        id: 12,
        player_id: '550e8400-e29b-41d4-a716-446655440005',
        game_id: 101,
        status: 'added',
        points: 3,
        started_at: '2024-02-01T00:00:00Z',
        completed_at: null,
        comment: null,
        rating: null,
        play_time: null
    },
    {
        id: 13,
        player_id: '550e8400-e29b-41d4-a716-446655440005',
        game_id: 103,
        status: 'completed',
        points: 1,
        started_at: '2024-01-25T00:00:00Z',
        completed_at: '2024-01-25T18:00:00Z',
        comment: 'Короткая, но хорошая',
        rating: 85,
        play_time: 'PT2H'
    },
    {
        id: 14,
        player_id: '550e8400-e29b-41d4-a716-446655440006',
        game_id: 105,
        status: 'in_progress',
        points: 4,
        started_at: '2024-01-29T00:00:00Z',
        completed_at: null,
        comment: 'Нравится',
        rating: null,
        play_time: 'PT10H'
    },
    {
        id: 15,
        player_id: '550e8400-e29b-41d4-a716-446655440007',
        game_id: 107,
        status: 'completed',
        points: 2,
        started_at: '2024-01-20T00:00:00Z',
        completed_at: '2024-01-21T00:00:00Z',
        comment: 'Отличная игра',
        rating: 89,
        play_time: 'PT7H15M'
    },
];
