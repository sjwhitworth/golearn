import numpy as np
import itertools

from sklearn import mixture
from sklearn import preprocessing

# Number of samples per component
n_samples = 500

# Generate random sample, two components
np.random.seed(0)
C = np.array([[0., -0.1], [1.7, .4]])
X = np.r_[np.dot(np.random.randn(n_samples, 2), C),
                  .7 * np.random.randn(n_samples, 2) + np.array([-6, 3])]

gmm = mixture.GaussianMixture(n_components=2, init_params='random')
gmm.fit(X)
labels = gmm.predict(X)

with open('gaussian_mixture.csv', 'w') as fout:
    with open('gaussian_mixture_labels.csv', 'w') as flabout:
        for i in range(1000):
            fout.write(",".join([str(x) for x in X[i,:]]) + "\n")
            flabout.write(str(labels[i]) + "\n")
