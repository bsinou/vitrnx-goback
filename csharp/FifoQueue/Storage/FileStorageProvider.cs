using System;
using System.Collections.Generic;
using System.IO;
using System.Text.Json;

namespace FifoQueue.Storage
{
    /// <summary>
    /// File-based storage provider that persists queue items as JSON.
    /// Items are loaded into memory on construction and flushed to disk on every write operation.
    /// </summary>
    /// <typeparam name="T">The type of elements stored. Must be serializable by System.Text.Json.</typeparam>
    public sealed class FileStorageProvider<T> : IStorageProvider<T>
    {
        private readonly string _filePath;
        private readonly LinkedList<T> _list;
        private readonly JsonSerializerOptions _jsonOptions;

        /// <summary>
        /// Creates a new file-based storage provider.
        /// </summary>
        /// <param name="filePath">The path to the file used for persistence.</param>
        /// <exception cref="ArgumentException">Thrown when <paramref name="filePath"/> is null or whitespace.</exception>
        public FileStorageProvider(string filePath)
        {
            if (string.IsNullOrWhiteSpace(filePath))
                throw new ArgumentException("File path must not be null or empty.", nameof(filePath));

            _filePath = filePath;
            _jsonOptions = new JsonSerializerOptions { WriteIndented = false };
            _list = LoadFromDisk();
        }

        /// <inheritdoc/>
        public void Append(T item)
        {
            _list.AddLast(item);
            FlushToDisk();
        }

        /// <inheritdoc/>
        public T RemoveFirst()
        {
            if (_list.First is null)
                throw new InvalidOperationException("The queue is empty.");

            var value = _list.First.Value;
            _list.RemoveFirst();
            FlushToDisk();
            return value;
        }

        /// <inheritdoc/>
        public bool TryRemoveFirst(out T? item)
        {
            if (_list.First is null)
            {
                item = default;
                return false;
            }

            item = _list.First.Value;
            _list.RemoveFirst();
            FlushToDisk();
            return true;
        }

        /// <inheritdoc/>
        public T PeekFirst()
        {
            if (_list.First is null)
                throw new InvalidOperationException("The queue is empty.");

            return _list.First.Value;
        }

        /// <inheritdoc/>
        public bool TryPeekFirst(out T? item)
        {
            if (_list.First is null)
            {
                item = default;
                return false;
            }

            item = _list.First.Value;
            return true;
        }

        /// <inheritdoc/>
        public int Count => _list.Count;

        /// <inheritdoc/>
        public void Clear()
        {
            _list.Clear();
            FlushToDisk();
        }

        /// <inheritdoc/>
        public void Dispose()
        {
            // Flush remaining state to ensure persistence.
            FlushToDisk();
        }

        private LinkedList<T> LoadFromDisk()
        {
            try
            {
                if (!System.IO.File.Exists(_filePath))
                    return new LinkedList<T>();

                var json = System.IO.File.ReadAllText(_filePath);
                if (string.IsNullOrWhiteSpace(json))
                    return new LinkedList<T>();

                var items = JsonSerializer.Deserialize<List<T>>(json, _jsonOptions);
                return items is null ? new LinkedList<T>() : new LinkedList<T>(items);
            }
            catch (FileNotFoundException)
            {
                return new LinkedList<T>();
            }
        }

        private void FlushToDisk()
        {
            var directory = Path.GetDirectoryName(_filePath);
            if (!string.IsNullOrEmpty(directory))
                Directory.CreateDirectory(directory);

            var items = new List<T>(_list);
            var json = JsonSerializer.Serialize(items, _jsonOptions);

            // Write to a temporary file first, then move atomically to prevent corruption.
            var tempPath = _filePath + ".tmp";
            try
            {
                System.IO.File.WriteAllText(tempPath, json);
                System.IO.File.Move(tempPath, _filePath, overwrite: true);
            }
            catch
            {
                try { System.IO.File.Delete(tempPath); } catch { /* best effort cleanup */ }
                throw;
            }
        }
    }
}
