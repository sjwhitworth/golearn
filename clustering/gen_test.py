#
# Generate sample data for the DBSCAN test 
# 
# Lifted from http://scikit-learn.org/stable/auto_examples/cluster/plot_dbscan.html#example-cluster-plot-dbscan-py
#

import numpy as np

from sklearn.cluster import DBSCAN
from sklearn import metrics
from sklearn.datasets.samples_generator import make_blobs
from sklearn.preprocessing import StandardScaler

centers = [[1, 1], [-1, -1], [1, -1]]
X, labels_true = make_blobs(n_samples=750, centers=centers, cluster_std=0.4,
    random_state=0)

X = StandardScaler().fit_transform(X)
X = X.astype(np.float64)
db = DBSCAN(eps=0.3, min_samples=10, metric='l2', algorithm='brute').fit(X)
core_samples_mask = np.zeros_like(db.labels_, dtype=bool)
core_samples_mask[db.core_sample_indices_] = True
labels = db.labels_

with open('dbscan.csv', 'w') as fscanout:
    with open('dbscan_labels.csv', 'w') as fscanlabout:
        for i in range(750):
            fscanout.write(",".join([str(x) for x in X[i,:]]) + "\n")
            fscanlabout.write(str(labels[i]) + "\n")

