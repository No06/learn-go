-- PostgreSQL schema & seed data for learn-go project testing
-- Run on a clean database. Adjust schema name or wrap in transaction as needed.

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Drop existing tables in dependency order (if they exist)
DROP TABLE IF EXISTS note_comments CASCADE;
DROP TABLE IF EXISTS notes CASCADE;
DROP TABLE IF EXISTS message_receipts CASCADE;
DROP TABLE IF EXISTS messages CASCADE;
DROP TABLE IF EXISTS conversation_members CASCADE;
DROP TABLE IF EXISTS conversations CASCADE;
DROP TABLE IF EXISTS submission_comments CASCADE;
DROP TABLE IF EXISTS submission_items CASCADE;
DROP TABLE IF EXISTS assignment_submissions CASCADE;
DROP TABLE IF EXISTS assignment_questions CASCADE;
DROP TABLE IF EXISTS assignments CASCADE;
DROP TABLE IF EXISTS course_sessions CASCADE;
DROP TABLE IF EXISTS courses CASCADE;
DROP TABLE IF EXISTS course_slots CASCADE;
DROP TABLE IF EXISTS teacher_student_links CASCADE;
DROP TABLE IF EXISTS students CASCADE;
DROP TABLE IF EXISTS teachers CASCADE;
DROP TABLE IF EXISTS classes CASCADE;
DROP TABLE IF EXISTS departments CASCADE;
DROP TABLE IF EXISTS accounts CASCADE;
DROP TABLE IF EXISTS schools CASCADE;

-- Core reference tables ------------------------------------------------------
CREATE TABLE schools (
    id          UUID PRIMARY KEY,
    name        VARCHAR(128) UNIQUE NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE accounts (
    id            UUID PRIMARY KEY,
    school_id     UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    role          VARCHAR(16) NOT NULL,
    identifier    VARCHAR(64) NOT NULL UNIQUE,
    password_hash VARCHAR(128) NOT NULL,
    display_name  VARCHAR(128) NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at    TIMESTAMPTZ
);

CREATE INDEX accounts_role_idx ON accounts(role);
CREATE INDEX accounts_school_idx ON accounts(school_id);

CREATE TABLE departments (
    id         UUID PRIMARY KEY,
    school_id  UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    name       VARCHAR(128) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE classes (
    id            UUID PRIMARY KEY,
    school_id     UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    department_id UUID NOT NULL REFERENCES departments(id) ON DELETE CASCADE,
    name          VARCHAR(128) NOT NULL,
    homeroom_id   UUID,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE teachers (
    id         UUID PRIMARY KEY,
    school_id  UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    account_id UUID NOT NULL UNIQUE REFERENCES accounts(id) ON DELETE CASCADE,
    number     VARCHAR(64) NOT NULL UNIQUE,
    email      VARCHAR(128) NOT NULL,
    phone      VARCHAR(32),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE students (
    id         UUID PRIMARY KEY,
    school_id  UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    account_id UUID NOT NULL UNIQUE REFERENCES accounts(id) ON DELETE CASCADE,
    number     VARCHAR(64) NOT NULL UNIQUE,
    class_id   UUID NOT NULL REFERENCES classes(id) ON DELETE CASCADE,
    email      VARCHAR(128) NOT NULL,
    phone      VARCHAR(32),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE teacher_student_links (
    id         UUID PRIMARY KEY,
    teacher_id UUID NOT NULL REFERENCES teachers(id) ON DELETE CASCADE,
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (teacher_id, student_id)
);

-- Course/assignment tables ---------------------------------------------------
CREATE TABLE courses (
    id          UUID PRIMARY KEY,
    school_id   UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    name        VARCHAR(128) NOT NULL,
    description VARCHAR(512),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE assignments (
    id             UUID PRIMARY KEY,
    course_id      UUID NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    teacher_id     UUID NOT NULL REFERENCES teachers(id) ON DELETE CASCADE,
    class_id       UUID NOT NULL REFERENCES classes(id) ON DELETE CASCADE,
    type           VARCHAR(16) NOT NULL,
    title          VARCHAR(256) NOT NULL,
    description    VARCHAR(1024),
    start_at       TIMESTAMPTZ,
    due_at         TIMESTAMPTZ,
    max_score      NUMERIC(6,2) NOT NULL,
    allow_resubmit BOOLEAN NOT NULL DEFAULT FALSE,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE assignment_questions (
    id            UUID PRIMARY KEY,
    assignment_id UUID NOT NULL REFERENCES assignments(id) ON DELETE CASCADE,
    type          VARCHAR(16) NOT NULL,
    prompt        TEXT NOT NULL,
    options       TEXT,
    answer        TEXT,
    score         NUMERIC(6,2) NOT NULL,
    order_index   INT NOT NULL
);

CREATE TABLE assignment_submissions (
    id            UUID PRIMARY KEY,
    assignment_id UUID NOT NULL REFERENCES assignments(id) ON DELETE CASCADE,
    student_id    UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    submitted_at  TIMESTAMPTZ,
    score         NUMERIC(6,2),
    feedback      TEXT,
    status        VARCHAR(32) NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX assignment_submissions_unique ON assignment_submissions(assignment_id, student_id);

CREATE TABLE submission_items (
    id            UUID PRIMARY KEY,
    submission_id UUID NOT NULL REFERENCES assignment_submissions(id) ON DELETE CASCADE,
    question_id   UUID NOT NULL REFERENCES assignment_questions(id) ON DELETE CASCADE,
    answer        TEXT,
    score         NUMERIC(6,2)
);

CREATE TABLE submission_comments (
    id            UUID PRIMARY KEY,
    submission_id UUID NOT NULL REFERENCES assignment_submissions(id) ON DELETE CASCADE,
    author_id     UUID NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
    author_role   VARCHAR(16) NOT NULL,
    content       TEXT NOT NULL,
    attachment_uri VARCHAR(256),
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Conversation/messaging tables ---------------------------------------------
CREATE TABLE conversations (
    id         UUID PRIMARY KEY,
    type       VARCHAR(16) NOT NULL,
    school_id  UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE conversation_members (
    id              UUID PRIMARY KEY,
    conversation_id UUID NOT NULL REFERENCES conversations(id) ON DELETE CASCADE,
    account_id      UUID NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
    role            VARCHAR(16) NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (conversation_id, account_id)
);

CREATE TABLE messages (
    id              UUID PRIMARY KEY,
    conversation_id UUID NOT NULL REFERENCES conversations(id) ON DELETE CASCADE,
    sender_id       UUID NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
    sender_role     VARCHAR(16) NOT NULL,
    kind            VARCHAR(16) NOT NULL,
    text            TEXT,
    media_uri       VARCHAR(256),
    metadata        TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE message_receipts (
    id         UUID PRIMARY KEY,
    message_id UUID NOT NULL REFERENCES messages(id) ON DELETE CASCADE,
    account_id UUID NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
    read_at    TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (message_id, account_id)
);

-- Notes ---------------------------------------------------------------------
CREATE TABLE notes (
    id         UUID PRIMARY KEY,
    school_id  UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    owner_id   UUID NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
    owner_role VARCHAR(16) NOT NULL,
    title      VARCHAR(256) NOT NULL,
    content    TEXT NOT NULL,
    visibility VARCHAR(16) NOT NULL,
    status     VARCHAR(16) NOT NULL,
    deleted_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE note_comments (
    id         UUID PRIMARY KEY,
    note_id    UUID NOT NULL REFERENCES notes(id) ON DELETE CASCADE,
    author_id  UUID NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
    author_role VARCHAR(16) NOT NULL,
    content    TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Seed data -----------------------------------------------------------------
INSERT INTO schools (id, name) VALUES
    ('11111111-1111-1111-1111-111111111111', 'Horizon International School');

INSERT INTO accounts (id, school_id, role, identifier, password_hash, display_name)
VALUES
    ('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', '11111111-1111-1111-1111-111111111111', 'admin',   'admin001',   '$2a$10$examplehashedpasswordadmin', '校区管理员'),
    ('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', '11111111-1111-1111-1111-111111111111', 'teacher', 'tch-1001',   '$2a$10$examplehashedpasswordteach', '李老师'),
    ('cccccccc-cccc-cccc-cccc-cccccccccccc', '11111111-1111-1111-1111-111111111111', 'student', 'stu-2025001', '$2a$10$examplehashedpasswordstud', '张三');

INSERT INTO departments (id, school_id, name) VALUES
    ('22222222-2222-2222-2222-222222222222', '11111111-1111-1111-1111-111111111111', '信息工程系');

INSERT INTO classes (id, school_id, department_id, name, homeroom_id)
VALUES ('33333333-3333-3333-3333-333333333333', '11111111-1111-1111-1111-111111111111', '22222222-2222-2222-2222-222222222222', '信工 2025 级 1 班', NULL);

INSERT INTO teachers (id, school_id, account_id, number, email, phone)
VALUES ('44444444-4444-4444-4444-444444444444', '11111111-1111-1111-1111-111111111111', 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 'tch-1001', 'teacher@example.com', '13800001111');

INSERT INTO students (id, school_id, account_id, number, class_id, email, phone)
VALUES ('55555555-5555-5555-5555-555555555555', '11111111-1111-1111-1111-111111111111', 'cccccccc-cccc-cccc-cccc-cccccccccccc', 'stu-2025001', '33333333-3333-3333-3333-333333333333', 'student@example.com', '13900002222');

INSERT INTO teacher_student_links (id, teacher_id, student_id)
VALUES ('66666666-6666-6666-6666-666666666666', '44444444-4444-4444-4444-444444444444', '55555555-5555-5555-5555-555555555555');

INSERT INTO courses (id, school_id, name, description)
VALUES ('77777777-7777-7777-7777-777777777777', '11111111-1111-1111-1111-111111111111', '计算机网络', '大二核心课程');

INSERT INTO assignments (id, course_id, teacher_id, class_id, type, title, description, start_at, due_at, max_score, allow_resubmit)
VALUES (
    '88888888-8888-8888-8888-888888888888',
    '77777777-7777-7777-7777-777777777777',
    '44444444-4444-4444-4444-444444444444',
    '33333333-3333-3333-3333-333333333333',
    'homework',
    '第 1 章作业',
    '阅读教材第 1 章并回答问题',
    NOW() - INTERVAL '7 days',
    NOW() + INTERVAL '3 days',
    100,
    TRUE
);

INSERT INTO assignment_questions (id, assignment_id, type, prompt, options, answer, score, order_index)
VALUES
    ('99999999-9999-9999-9999-999999999999', '88888888-8888-8888-8888-888888888888', 'essay', '解释 OSI 七层模型的每一层职责。', NULL, '参考教材示例答案', 40, 1),
    ('aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee', '88888888-8888-8888-8888-888888888888', 'choice', '第 3 层是以下哪一项？', '{"options":["会话层","网络层","物理层"]}', '网络层', 60, 2);

INSERT INTO assignment_submissions (id, assignment_id, student_id, submitted_at, score, feedback, status)
VALUES ('bbbbbbbb-cccc-dddd-eeee-ffffffffffff', '88888888-8888-8888-8888-888888888888', '55555555-5555-5555-5555-555555555555', NOW() - INTERVAL '2 days', 88.5, '整体表现良好，复习第 3 题', 'graded');

INSERT INTO submission_items (id, submission_id, question_id, answer, score)
VALUES
    ('cccccccc-dddd-eeee-ffff-000000000000', 'bbbbbbbb-cccc-dddd-eeee-ffffffffffff', '99999999-9999-9999-9999-999999999999', '从应用层到物理层逐一说明', 38),
    ('dddddddd-eeee-ffff-0000-111111111111', 'bbbbbbbb-cccc-dddd-eeee-ffffffffffff', 'aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee', '选择：网络层', 50.5);

INSERT INTO submission_comments (id, submission_id, author_id, author_role, content)
VALUES ('eeeeeeee-ffff-0000-1111-222222222222', 'bbbbbbbb-cccc-dddd-eeee-ffffffffffff', 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 'teacher', '请再次阅读教材第 2 章，加深理解。');

INSERT INTO conversations (id, type, school_id)
VALUES ('ffffffff-0000-1111-2222-333333333333', 'direct', '11111111-1111-1111-1111-111111111111');

INSERT INTO conversation_members (id, conversation_id, account_id, role)
VALUES
    ('00000000-1111-2222-3333-444444444444', 'ffffffff-0000-1111-2222-333333333333', 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 'teacher'),
    ('11111111-2222-3333-4444-555555555555', 'ffffffff-0000-1111-2222-333333333333', 'cccccccc-cccc-cccc-cccc-cccccccccccc', 'student');

INSERT INTO messages (id, conversation_id, sender_id, sender_role, kind, text)
VALUES ('22222222-3333-4444-5555-666666666666', 'ffffffff-0000-1111-2222-333333333333', 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 'teacher', 'text', '记得周五前提交实验报告。');

INSERT INTO message_receipts (id, message_id, account_id, read_at)
VALUES ('33333333-4444-5555-6666-777777777777', '22222222-3333-4444-5555-666666666666', 'cccccccc-cccc-cccc-cccc-cccccccccccc', NOW());

INSERT INTO notes (id, school_id, owner_id, owner_role, title, content, visibility, status)
VALUES ('44444444-5555-6666-7777-888888888888', '11111111-1111-1111-1111-111111111111', 'cccccccc-cccc-cccc-cccc-cccccccccccc', 'student', '网络实验心得', '记录一次路由实验的心得体会。', 'class', 'published');

INSERT INTO note_comments (id, note_id, author_id, author_role, content)
VALUES ('55555555-6666-7777-8888-999999999999', '44444444-5555-6666-7777-888888888888', 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 'teacher', '很好，补充下实验截图会更完整。');
