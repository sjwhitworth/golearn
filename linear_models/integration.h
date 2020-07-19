#ifndef _H_INTEGRATION_
#define _H_INTEGRATION_

#include "linear.h"

struct problem *CreateCProblem();
void FreeCProblem(struct problem*);
struct model *CreateCModel();
void FreeCModel(struct model*);
struct parameter *CreateCParameter();
void FreeCParameter(struct parameter*);
// Allocates memory outside of golang for describing feature
// vectors.
int RiffleFeatures(struct problem *p, int num_offsets, int* row_offsets, struct feature_node *features);
int AllocateLabelsForProblem(struct problem *, int);
void AssignLabelForProblem(struct problem *, int, double);
struct feature_node *GetFeatureNodeForIndex(struct problem *, int, int);

#endif
