using System;
using System.Collections.Generic;

namespace FifoQueue.Storage
{
    /// <summary>
    /// In-memory storage provider backed by a <see cref="LinkedList{T}"/>.
    /// Provides O(1) enqueue and dequeue operations.
    /// </summary>
    /// <typeparam name="T">The type of elements stored.</typeparam>
    public sealed class MemoryStorageProvider<T> : IStorageProvider<T>
    {
        private readonly LinkedList<T> _list = new();

        /// <inheritdoc/>
        public void Append(T item) => _list.AddLast(item);

        /// <inheritdoc/>
        public T RemoveFirst()
        {
            if (_list.First is null)
                throw new InvalidOperationException("The queue is empty.");

            var value = _list.First.Value;
            _list.RemoveFirst();
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
        public void Clear() => _list.Clear();

        /// <inheritdoc/>
        public void Dispose()
        {
            // No unmanaged resources to release.
            _list.Clear();
        }
    }
}
