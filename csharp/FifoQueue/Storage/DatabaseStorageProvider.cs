using System;
using Microsoft.Data.Sqlite;

namespace FifoQueue.Storage
{
    /// <summary>
    /// SQLite-based storage provider that persists queue items in a database table.
    /// Items are serialized to JSON for storage.
    /// </summary>
    /// <typeparam name="T">The type of elements stored. Must be serializable by System.Text.Json.</typeparam>
    public sealed class DatabaseStorageProvider<T> : IStorageProvider<T>
    {
        private readonly SqliteConnection _connection;

        /// <summary>
        /// Creates a new database storage provider.
        /// </summary>
        /// <param name="connectionString">
        /// SQLite connection string.
        /// Examples: "Data Source=queue.db" for a file-based database,
        /// or "Data Source=:memory:" for an in-memory database.
        /// </param>
        /// <exception cref="ArgumentException">Thrown when <paramref name="connectionString"/> is null or whitespace.</exception>
        public DatabaseStorageProvider(string connectionString)
        {
            if (string.IsNullOrWhiteSpace(connectionString))
                throw new ArgumentException("Connection string must not be null or empty.", nameof(connectionString));

            _connection = new SqliteConnection(connectionString);
            _connection.Open();
            InitializeSchema();
        }

        /// <inheritdoc/>
        public void Append(T item)
        {
            var json = System.Text.Json.JsonSerializer.Serialize(item);
            using var cmd = _connection.CreateCommand();
            cmd.CommandText = "INSERT INTO queue_items (payload) VALUES (@payload);";
            cmd.Parameters.AddWithValue("@payload", json);
            cmd.ExecuteNonQuery();
        }

        /// <inheritdoc/>
        public T RemoveFirst()
        {
            using var cmd = _connection.CreateCommand();
            cmd.CommandText = "DELETE FROM queue_items WHERE id = (SELECT MIN(id) FROM queue_items) RETURNING payload;";
            var result = cmd.ExecuteScalar();

            if (result is null or DBNull)
                throw new InvalidOperationException("The queue is empty.");

            return System.Text.Json.JsonSerializer.Deserialize<T>((string)result)!;
        }

        /// <inheritdoc/>
        public bool TryRemoveFirst(out T? item)
        {
            using var cmd = _connection.CreateCommand();
            cmd.CommandText = "DELETE FROM queue_items WHERE id = (SELECT MIN(id) FROM queue_items) RETURNING payload;";
            var result = cmd.ExecuteScalar();

            if (result is null or DBNull)
            {
                item = default;
                return false;
            }

            item = System.Text.Json.JsonSerializer.Deserialize<T>((string)result);
            return true;
        }

        /// <inheritdoc/>
        public T PeekFirst()
        {
            using var cmd = _connection.CreateCommand();
            cmd.CommandText = "SELECT payload FROM queue_items ORDER BY id ASC LIMIT 1;";
            var result = cmd.ExecuteScalar();

            if (result is null or DBNull)
                throw new InvalidOperationException("The queue is empty.");

            return System.Text.Json.JsonSerializer.Deserialize<T>((string)result)!;
        }

        /// <inheritdoc/>
        public bool TryPeekFirst(out T? item)
        {
            using var cmd = _connection.CreateCommand();
            cmd.CommandText = "SELECT payload FROM queue_items ORDER BY id ASC LIMIT 1;";
            var result = cmd.ExecuteScalar();

            if (result is null or DBNull)
            {
                item = default;
                return false;
            }

            item = System.Text.Json.JsonSerializer.Deserialize<T>((string)result);
            return true;
        }

        /// <inheritdoc/>
        public int Count
        {
            get
            {
                using var cmd = _connection.CreateCommand();
                cmd.CommandText = "SELECT COUNT(*) FROM queue_items;";
                return Convert.ToInt32(cmd.ExecuteScalar());
            }
        }

        /// <inheritdoc/>
        public void Clear()
        {
            using var cmd = _connection.CreateCommand();
            cmd.CommandText = "DELETE FROM queue_items;";
            cmd.ExecuteNonQuery();
        }

        /// <inheritdoc/>
        public void Dispose()
        {
            _connection.Dispose();
        }

        private void InitializeSchema()
        {
            using var cmd = _connection.CreateCommand();
            cmd.CommandText = @"
                CREATE TABLE IF NOT EXISTS queue_items (
                    id INTEGER PRIMARY KEY AUTOINCREMENT,
                    payload TEXT NOT NULL
                );";
            cmd.ExecuteNonQuery();
        }
    }
}
