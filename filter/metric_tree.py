import numpy as np

class Tree:
    root = None
    tolerance = 3

    class Node(object):
        def __init__(self, word):
            self.word = word
            self.children = {}

    def insert_word(self, word):
        if self.root == None:
            self.root = self.Node(word)
        else:
            self.__insert_word_internal(self.Node(word))

    def __insert_word_internal(self, new_node):
        parent = self.root

        dist = self.__distance(new_node.word, parent.word)

        while dist in parent.children:
            parent = parent.children[dist]
            dist = self.__distance(new_node.word, parent.word)
        
        parent.children[dist] = new_node

    def get_score(self, word, modifier):
        results = []
        corrections = self.__get_score_rec(word, self.root)
        for node in corrections:
            results.append(node.word)
        dist_list = [self.__distance(results[i], word) for i in range(len(results))]

        if len(dist_list) > 0:
            closest = min(dist_list)
            score = (len(word) - closest + modifier)/len(word)
        else:
            score = 0

        return score

    def __get_score_rec(self, word, root):
        dist = self.__distance(word, root.word)
        similar_words = []
        if dist < self.tolerance:
            similar_words.append(root)

        min_val = max(1, dist - self.tolerance)
        max_val = dist + self.tolerance

        for i in range(min_val, max_val+1):
            if i in root.children:
                similar_words += self.__get_score_rec(word, root.children[i])

        return similar_words

    def __distance(self, s1, s2):
        if (len(s1) == 0):
            return len(s2)
        if (len(s2) == 0):
            return len(s1)
        dist = np.zeros(shape=(len(s1)+1, len(s2)+1))

        for i in range(1, len(s1)+1):
            dist[i][0] = i

            for j in range(1, len(s2)+1):
                cost = not s1[i - 1] == s2[j - 1]
                if (i == 1):
                    dist[0][j] = j

                dist[i][j] = self.self_min( \
                    dist[i - 1][j] + 1, \
                    [dist[i][j - 1] + 1, \
                    dist[i - 1][j - 1] + cost])
                if (i > 1 and j > 1 and s1[i - 1] == s2[j - 2] and s1[i - 2] == s2[j - 1]):
                    dist[i][j] = min(dist[i][j], dist[i - 2][j - 2] + cost)
        
        return int(dist[len(s1)][len(s2)])

    def self_min(self, s1, vals):
        if len(vals) == 0:
            return s1
        current_min = s1
        for i in vals:
            current_min = min(current_min, i)
        return current_min

# tree = Tree()

# # tree.insert_word("test")
# tree.insert_word("testt")
# tree.insert_word("tst")

# tree.insert_word("Test")
# tree.insert_word("tesT")


# # print(tree.root.children[1].word)
# print(tree.get_score("test"))


# # print(tree.default_tolerance)