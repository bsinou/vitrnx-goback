using System;

namespace FifoQueue
{
    /// <summary>
    /// A thread-safe FIFO queue that delegates storage to an <see cref="IStorageProvider{T}"/>.
    /// All public operations are protected by a lock for concurrent access safety.
    /// </summary>
    /// <typeparam name="T">The type of elements in the queue.</typeparam>
    public sealed class FifoQueue<T> : IFifoQueue<T>
    {
        private readonly IStorageProvider<T> _storage;
        private readonly object _lock = new();
        private bool _disposed;

        /// <summary>
        /// Creates a new FIFO queue backed by the specified storage provider.
        /// </summary>
        /// <param name="storageProvider">The storage provider to use.</param>
        /// <exception cref="ArgumentNullException">Thrown when <paramref name="storageProvider"/> is null.</exception>
        public FifoQueue(IStorageProvider<T> storageProvider)
        {
            _storage = storageProvider ?? throw new ArgumentNullException(nameof(storageProvider));
        }

        /// <inheritdoc/>
        public void Enqueue(T item)
        {
            lock (_lock)
            {
                ThrowIfDisposed();
                _storage.Append(item);
            }
        }

        /// <inheritdoc/>
        public T Dequeue()
        {
            lock (_lock)
            {
                ThrowIfDisposed();
                return _storage.RemoveFirst();
            }
        }

        /// <inheritdoc/>
        public bool TryDequeue(out T? item)
        {
            lock (_lock)
            {
                ThrowIfDisposed();
                return _storage.TryRemoveFirst(out item);
            }
        }

        /// <inheritdoc/>
        public T Peek()
        {
            lock (_lock)
            {
                ThrowIfDisposed();
                return _storage.PeekFirst();
            }
        }

        /// <inheritdoc/>
        public bool TryPeek(out T? item)
        {
            lock (_lock)
            {
                ThrowIfDisposed();
                return _storage.TryPeekFirst(out item);
            }
        }

        /// <inheritdoc/>
        public int Count
        {
            get
            {
                lock (_lock)
                {
                    ThrowIfDisposed();
                    return _storage.Count;
                }
            }
        }

        /// <inheritdoc/>
        public bool IsEmpty
        {
            get
            {
                lock (_lock)
                {
                    ThrowIfDisposed();
                    return _storage.Count == 0;
                }
            }
        }

        /// <inheritdoc/>
        public void Clear()
        {
            lock (_lock)
            {
                ThrowIfDisposed();
                _storage.Clear();
            }
        }

        /// <inheritdoc/>
        public void Dispose()
        {
            lock (_lock)
            {
                if (_disposed)
                    return;

                _disposed = true;
                _storage.Dispose();
            }
        }

        private void ThrowIfDisposed()
        {
            if (_disposed)
                throw new ObjectDisposedException(nameof(FifoQueue<T>));
        }
    }
}
